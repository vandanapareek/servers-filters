package services

import (
	"context"

	"servers-filters/dto"
)

// interface for server business logic
type ServerService interface {
	GetServers(ctx context.Context, req dto.ServerListRequest) (*dto.ServerListResponse, error)
	GetLocations(ctx context.Context) ([]string, error)
	GetMetrics(ctx context.Context) (*dto.MetricsResponse, error)
}
