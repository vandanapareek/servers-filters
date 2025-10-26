package repository

import (
	"context"

	"servers-filters/models"
)

//  Server data operations interface
type ServerRepository interface {
	GetServers(ctx context.Context, filters models.ServerFilters) ([]models.Server, int64, error)

	GetServerCount(ctx context.Context, filters models.ServerFilters) (int64, error)

	GetLocations(ctx context.Context) ([]string, error)

	GetMetrics(ctx context.Context) (*models.ServerMetrics, error)
}
