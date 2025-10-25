package services

import (
	"context"

	"servers-filters/dto"
)

// ServerService defines the interface for server business logic
type ServerService interface {
	// Get servers with filters and pagination
	GetServers(ctx context.Context, req dto.ServerListRequest) (*dto.ServerListResponse, error)


	// Get all unique locations
	GetLocations(ctx context.Context) ([]string, error)

	// Get metrics about the servers
	GetMetrics(ctx context.Context) (*dto.MetricsResponse, error)
}
