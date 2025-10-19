package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"servers-filters/handlers"
	"servers-filters/internal/config"
	"servers-filters/internal/logger"
	"servers-filters/repository"
	"servers-filters/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	log := logger.GetLogger()

	log.Info("Starting servers listing application")

	// Initialize database
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}
	defer db.Close()

	// Initialize repositories
	serverRepo := repository.NewSQLiteRepository(db)

	var cacheRepo repository.CacheRepository
	if cfg.Cache.Enabled {
		redisClient, err := initRedis(cfg.Cache)
		if err != nil {
			log.WithError(err).Warn("Failed to initialize Redis, continuing without cache")
			cacheRepo = repository.NewNoOpCacheRepository()
		} else {
			cacheRepo = repository.NewRedisCacheRepository(redisClient)
			log.Info("Redis cache initialized")
		}
	} else {
		cacheRepo = repository.NewNoOpCacheRepository()
		log.Info("Cache disabled")
	}

	// Initialize services
	serverService := services.NewServerService(serverRepo, cacheRepo, cfg.Cache.TTL)

	// Initialize handlers
	serverHandler := handlers.NewServerHandler(serverService)

	// Setup router
	router := setupRouter(serverHandler)

	// Create server
	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("addr", server.Addr).Info("Server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server exited")
}

// initDatabase initializes the database connection
func initDatabase(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// initRedis initializes the Redis connection
func initRedis(cfg config.CacheConfig) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}

// setupRouter sets up the HTTP router with middleware and routes
func setupRouter(serverHandler *handlers.ServerHandler) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Server routes
		r.Get("/servers", serverHandler.GetServers)
		r.Get("/servers/{id}", serverHandler.GetServerByID)
		r.Get("/locations", serverHandler.GetLocations)
		r.Get("/metrics", serverHandler.GetMetrics)
	})

	// Legacy routes for backward compatibility
	router.Get("/servers", serverHandler.GetServers)
	router.Get("/servers/{id}", serverHandler.GetServerByID)
	router.Get("/locations", serverHandler.GetLocations)
	router.Get("/metrics", serverHandler.GetMetrics)

	return router
}
