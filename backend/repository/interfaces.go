package repository

import (
	"context"

	"servers-filters/models"
)

// Interface for server data operations
type ServerRepository interface {
	// Get servers with filters and pagination
	GetServers(ctx context.Context, filters models.ServerFilters) ([]models.Server, int64, error)

	// Get server by its ID
	GetServerByID(ctx context.Context, id int) (*models.Server, error)

	// Get total count of servers matching the filters
	GetServerCount(ctx context.Context, filters models.ServerFilters) (int64, error)

	// Get all unique locations
	GetLocations(ctx context.Context) ([]string, error)

	// Get metrics about the servers
	GetMetrics(ctx context.Context) (*models.ServerMetrics, error)
}

// Interface for cache operations
type CacheRepository interface {
	// Get a value from cache
	Get(ctx context.Context, key string) (string, error)

	// Set value in cache with TTL
	Set(ctx context.Context, key string, value string, ttl int) error

	// Delete value from cache
	Delete(ctx context.Context, key string) error

	// Clear all values from cache
	Clear(ctx context.Context) error
}
