package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"servers-filters/handlers"
	"servers-filters/internal/logger"
	"servers-filters/repository"
	"servers-filters/services"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// Create a temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Copy the test database
	sourceDB := "data/servers.db"
	if _, err := os.Stat(sourceDB); os.IsNotExist(err) {
		t.Skip("Test database not found, skipping integration test")
	}

	// Read and copy the database
	data, err := os.ReadFile(sourceDB)
	if err != nil {
		t.Fatalf("Failed to read source database: %v", err)
	}

	err = os.WriteFile(dbPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test database: %v", err)
	}

	// Connect to the test database
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

func setupTestServer(t *testing.T) *httptest.Server {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Initialize logger
	logger.Init("debug", "text")

	// Initialize repositories
	serverRepo := repository.NewSQLiteRepository(db)
	cacheRepo := repository.NewNoOpCacheRepository()

	// Initialize services
	serverService := services.NewServerService(serverRepo, cacheRepo, 300)

	// Initialize handlers
	serverHandler := handlers.NewServerHandler(serverService)

	// Create router
	router := setupTestRouter(serverHandler)

	return httptest.NewServer(router)
}

func setupTestRouter(serverHandler *handlers.ServerHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/servers", serverHandler.GetServers)
	mux.HandleFunc("/servers/", func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path
		path := r.URL.Path
		if len(path) > 8 { // "/servers/" = 8 chars
			id := path[8:]
			r.URL.Path = "/servers/" + id
		}
		serverHandler.GetServerByID(w, r)
	})
	mux.HandleFunc("/locations", serverHandler.GetLocations)
	mux.HandleFunc("/metrics", serverHandler.GetMetrics)
	mux.HandleFunc("/health", serverHandler.Health)

	return mux
}

func TestIntegration_GetServers(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Get all servers",
			url:            "/servers",
			expectedStatus: http.StatusOK,
			expectedCount:  5, // Based on our test data
		},
		{
			name:           "Get servers with pagination",
			url:            "/servers?page=1&per_page=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get servers with RAM filter",
			url:            "/servers?ram_min=32",
			expectedStatus: http.StatusOK,
			expectedCount:  2, // Servers with 32GB+ RAM
		},
		{
			name:           "Get servers with price filter",
			url:            "/servers?price_max=100",
			expectedStatus: http.StatusOK,
			expectedCount:  4, // Servers under €100
		},
		{
			name:           "Get servers with location filter",
			url:            "/servers?location=Amsterdam",
			expectedStatus: http.StatusOK,
			expectedCount:  2, // Servers in Amsterdam
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.url)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			data, ok := result["data"].([]interface{})
			if !ok {
				t.Fatalf("Response data is not an array")
			}

			if len(data) != tt.expectedCount {
				t.Errorf("Expected %d servers, got %d", tt.expectedCount, len(data))
			}
		})
	}
}

func TestIntegration_GetServerByID(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		id             int
		expectedStatus int
		shouldExist    bool
	}{
		{
			name:           "Get existing server",
			id:             1,
			expectedStatus: http.StatusOK,
			shouldExist:    true,
		},
		{
			name:           "Get non-existing server",
			id:             999,
			expectedStatus: http.StatusNotFound,
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/servers/%d", server.URL, tt.id))
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.shouldExist {
				var result map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				data, ok := result["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Response data is not an object")
				}

				if int(data["id"].(float64)) != tt.id {
					t.Errorf("Expected server ID %d, got %v", tt.id, data["id"])
				}
			}
		})
	}
}

func TestIntegration_GetLocations(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/locations")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatalf("Response data is not an array")
	}

	if len(data) == 0 {
		t.Error("Expected at least one location")
	}
}

func TestIntegration_GetMetrics(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check required fields
	requiredFields := []string{"total_servers", "average_price", "locations_count", "last_updated"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Check that total_servers is greater than 0
	if totalServers, ok := result["total_servers"].(float64); ok {
		if totalServers <= 0 {
			t.Error("Expected total_servers to be greater than 0")
		}
	}
}

func TestIntegration_HealthCheck(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if status, ok := result["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", result["status"])
	}
}
