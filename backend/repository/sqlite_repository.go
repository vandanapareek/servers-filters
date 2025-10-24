package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"servers-filters/models"

	"github.com/jmoiron/sqlx"
)

// SQLiteRepository implements ServerRepository for SQLite
type SQLiteRepository struct {
	db *sqlx.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(db *sqlx.DB) ServerRepository {
	return &SQLiteRepository{db: db}
}

// GetServers retrieves servers with filters and pagination
func (r *SQLiteRepository) GetServers(ctx context.Context, filters models.ServerFilters) ([]models.Server, int64, error) {
	// Build WHERE clause
	whereClause, args := r.buildWhereClause(filters)

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(filters.Sort)

	// Build LIMIT and OFFSET
	limit := filters.PerPage
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Build query
	query := fmt.Sprintf(`
		SELECT id, model, cpu, ram_gb, hdd, storage_gb, location_city, 
		       location_code, price_eur, raw_price, raw_ram, raw_hdd, created_at
		FROM servers
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	// Add limit and offset to args
	args = append(args, limit, offset)

	// Execute query
	var servers []models.Server
	fmt.Printf("DEBUG REPO: Query=%s, Args=%v\n", query, args)
	err := r.db.SelectContext(ctx, &servers, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get servers: %w", err)
	}
	fmt.Printf("DEBUG REPO: Returned %d servers\n", len(servers))

	// Also test the query manually
	fmt.Printf("DEBUG REPO: Testing query manually...\n")
	var testServers []models.Server
	testErr := r.db.SelectContext(ctx, &testServers, query, args...)
	if testErr != nil {
		fmt.Printf("DEBUG REPO: Manual test failed: %v\n", testErr)
	} else {
		fmt.Printf("DEBUG REPO: Manual test returned %d servers\n", len(testServers))
	}

	// Get total count
	total, err := r.GetServerCount(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get server count: %w", err)
	}

	return servers, total, nil
}

// GetServerByID retrieves a server by its ID
func (r *SQLiteRepository) GetServerByID(ctx context.Context, id int) (*models.Server, error) {
	query := `
		SELECT id, model, cpu, ram_gb, hdd, storage_gb, location_city, 
		       location_code, price_eur, raw_price, raw_ram, raw_hdd, created_at
		FROM servers
		WHERE id = ?
	`

	var server models.Server
	err := r.db.GetContext(ctx, &server, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get server by ID: %w", err)
	}

	return &server, nil
}

// GetServerCount returns the total count of servers matching the filters
func (r *SQLiteRepository) GetServerCount(ctx context.Context, filters models.ServerFilters) (int64, error) {
	whereClause, args := r.buildWhereClause(filters)

	query := fmt.Sprintf("SELECT COUNT(*) FROM servers %s", whereClause)

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get server count: %w", err)
	}

	return count, nil
}

// GetLocations returns all unique locations
func (r *SQLiteRepository) GetLocations(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT location_city 
		FROM servers 
		WHERE location_city IS NOT NULL AND location_city != ''
		ORDER BY location_city
	`

	var locations []string
	err := r.db.SelectContext(ctx, &locations, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	return locations, nil
}

// GetMetrics returns basic metrics about the servers
func (r *SQLiteRepository) GetMetrics(ctx context.Context) (*models.ServerMetrics, error) {
	query := `
		SELECT 
			COUNT(*) as total_servers,
			MIN(price_eur) as min_price,
			MAX(price_eur) as max_price,
			COUNT(DISTINCT location_city) as locations_count
		FROM servers
		WHERE price_eur IS NOT NULL AND location_city IS NOT NULL AND location_city != ''
	`

	var result struct {
		TotalServers   int64   `db:"total_servers"`
		MinPrice       float64 `db:"min_price"`
		MaxPrice       float64 `db:"max_price"`
		LocationsCount int64   `db:"locations_count"`
	}

	err := r.db.GetContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	metrics := &models.ServerMetrics{
		TotalServers:   result.TotalServers,
		MinPrice:       result.MinPrice,
		MaxPrice:       result.MaxPrice,
		LocationsCount: result.LocationsCount,
		LastUpdated:    time.Now(),
	}

	return metrics, nil
}

// buildWhereClause builds the WHERE clause and arguments for the query
func (r *SQLiteRepository) buildWhereClause(filters models.ServerFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// Text search
	if filters.Query != "" {
		conditions = append(conditions, "(model LIKE ? OR cpu LIKE ?)")
		searchTerm := "%" + filters.Query + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Location filter
	if len(filters.Location) > 0 {
		placeholders := make([]string, len(filters.Location))
		for i, loc := range filters.Location {
			placeholders[i] = "?"
			args = append(args, loc)
		}
		// Add the same values again for the second IN clause
		for _, loc := range filters.Location {
			args = append(args, loc)
		}
		conditions = append(conditions, fmt.Sprintf("(location_city IN (%s) OR location_code IN (%s))",
			strings.Join(placeholders, ","), strings.Join(placeholders, ",")))
	}

	// RAM exact values filter (for checkboxes)
	if len(filters.RAMValues) > 0 {
		placeholders := make([]string, len(filters.RAMValues))
		for i, ram := range filters.RAMValues {
			placeholders[i] = "?"
			args = append(args, ram)
		}
		conditions = append(conditions, fmt.Sprintf("ram_gb IN (%s)", strings.Join(placeholders, ",")))
	} else {
		// RAM range filter (fallback for min/max)
		if filters.RAMMin != nil {
			conditions = append(conditions, "ram_gb >= ?")
			args = append(args, *filters.RAMMin)
		}
		if filters.RAMMax != nil {
			conditions = append(conditions, "ram_gb <= ?")
			args = append(args, *filters.RAMMax)
		}
	}

	// Storage range filter (now receives GB values from service layer)
	if filters.StorageMin != nil {
		conditions = append(conditions, "storage_gb >= ?")
		args = append(args, *filters.StorageMin)
		fmt.Printf("🔥 CUSTOM DEBUG: StorageMin filter added: %d GB\n", *filters.StorageMin)
	}
	if filters.StorageMax != nil {
		conditions = append(conditions, "storage_gb <= ?")
		args = append(args, *filters.StorageMax)
		fmt.Printf("🔥 CUSTOM DEBUG: StorageMax filter added: %d GB\n", *filters.StorageMax)
	}

	// HDD filter
	if filters.HDD != "" {
		conditions = append(conditions, "hdd LIKE ?")
		args = append(args, "%"+filters.HDD+"%")
	}

	// Price range filter
	if filters.PriceMin != nil {
		conditions = append(conditions, "price_eur >= ?")
		args = append(args, *filters.PriceMin)
	}
	if filters.PriceMax != nil {
		conditions = append(conditions, "price_eur <= ?")
		args = append(args, *filters.PriceMax)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

// buildOrderClause builds the ORDER BY clause for the query
func (r *SQLiteRepository) buildOrderClause(sort string) string {
	sortOption := models.ParseSort(sort)
	return fmt.Sprintf("ORDER BY %s %s", sortOption.Field, strings.ToUpper(sortOption.Order))
}
