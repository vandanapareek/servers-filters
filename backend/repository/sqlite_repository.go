package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"servers-filters/internal/constants"
	"servers-filters/models"

	"github.com/jmoiron/sqlx"
)

// implement ServerRepository for SQLite
type SQLiteRepository struct {
	db *sqlx.DB
}

// create a new SQLite repository
func NewSQLiteRepository(db *sqlx.DB) ServerRepository {
	return &SQLiteRepository{db: db}
}

// get servers with filters and pagination
func (r *SQLiteRepository) GetServers(ctx context.Context, filters models.ServerFilters) ([]models.Server, int64, error) {
	whereClause, args := r.buildWhereClause(filters)
	orderClause := r.buildOrderClause(filters.Sort)
	limit := filters.PerPage
	if limit <= 0 {
		limit = constants.DefaultPerPage
	}
	if limit > constants.MaxPerPage {
		limit = constants.MaxPerPage
	}

	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Build query
	query := fmt.Sprintf(`
		SELECT id, model, cpu, ram_gb, hdd_gb, hdd_type, location, 
		       location_code, price, raw_price, raw_hdd, raw_ram, created_at, updated_at
		FROM servers
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	// Add limit and offset
	args = append(args, limit, offset)

	// Execute
	var servers []models.Server
	err := r.db.SelectContext(ctx, &servers, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get servers: %w", err)
	}

	// Get total count
	total, err := r.GetServerCount(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get server count: %w", err)
	}

	return servers, total, nil
}

// Get the total count of servers matching the filters
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

// Get all unique locations
func (r *SQLiteRepository) GetLocations(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT location 
		FROM servers 
		WHERE location IS NOT NULL AND location != ''
		ORDER BY location
	`

	var locations []string
	err := r.db.SelectContext(ctx, &locations, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	return locations, nil
}

// Get statistics about the servers
func (r *SQLiteRepository) GetMetrics(ctx context.Context) (*models.ServerMetrics, error) {
	query := `
		SELECT 
			COUNT(*) as total_servers,
			MIN(price) as min_price,
			MAX(price) as max_price,
			COUNT(DISTINCT location) as locations_count
		FROM servers
		WHERE price IS NOT NULL AND location IS NOT NULL AND location != ''
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

// build where clause and args for the query
func (r *SQLiteRepository) buildWhereClause(filters models.ServerFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// Model text search
	if filters.Query != "" {
		conditions = append(conditions, "model LIKE ?")
		searchTerm := "%" + filters.Query + "%"
		args = append(args, searchTerm)
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
		conditions = append(conditions, fmt.Sprintf("(location IN (%s) OR location_code IN (%s))",
			strings.Join(placeholders, ","), strings.Join(placeholders, ",")))
	}

	// RAM filter
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

	// Storage range filter (get GB values from service layer)
	if filters.StorageMin != nil {
		conditions = append(conditions, "hdd_gb >= ?")
		args = append(args, *filters.StorageMin)
	}
	if filters.StorageMax != nil {
		conditions = append(conditions, "hdd_gb <= ?")
		args = append(args, *filters.StorageMax)
	}

	// HDD type filter
	if filters.HDD != "" {
		conditions = append(conditions, "hdd_type = ?")
		args = append(args, filters.HDD)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

// build the order by clause for the query
func (r *SQLiteRepository) buildOrderClause(sort string) string {
	sortOption := models.ParseSort(sort)
	return fmt.Sprintf("ORDER BY %s %s", sortOption.Field, strings.ToUpper(sortOption.Order))
}
