package services

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"servers-filters/dto"
	"servers-filters/models"
	"servers-filters/repository"
)

// ServerServiceImpl implements ServerService
type ServerServiceImpl struct {
	serverRepo repository.ServerRepository
	cacheRepo  repository.CacheRepository
	cacheTTL   int
}

// Create a new server service
func NewServerService(serverRepo repository.ServerRepository, cacheRepo repository.CacheRepository, cacheTTL int) ServerService {
	return &ServerServiceImpl{
		serverRepo: serverRepo,
		cacheRepo:  cacheRepo,
		cacheTTL:   cacheTTL,
	}
}

// Get servers with filters and pagination
func (s *ServerServiceImpl) GetServers(ctx context.Context, req dto.ServerListRequest) (*dto.ServerListResponse, error) {
	// Normalize request parameters
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	// Convert request to model filters
	filters := s.convertRequestToFilters(req)

	// Generate cache key
	cacheKey := s.generateCacheKey("servers", filters)

	// Try to get from cache first
	var response *dto.ServerListResponse
	if s.cacheRepo != nil {
		cached, err := s.cacheRepo.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			err = json.Unmarshal([]byte(cached), &response)
			if err == nil {
				return response, nil
			}
		}
	}

	// Get from database
	servers, total, err := s.serverRepo.GetServers(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}

	// Apply pagination to the results
	start := (req.Page - 1) * req.PerPage
	end := start + req.PerPage
	
	if start >= len(servers) {
		servers = []models.Server{} // Empty slice if start is beyond available data
	} else if end > len(servers) {
		servers = servers[start:] // Take from start to end of available data
	} else {
		servers = servers[start:end] // Take the requested page slice
	}

	// Convert to DTOs
	serverDTOs := make([]dto.ServerDTO, len(servers))
	for i, server := range servers {
		serverDTOs[i] = s.convertModelToDTO(server)
	}

	// Calculate pagination
	totalPages := int((total + int64(req.PerPage) - 1) / int64(req.PerPage))

	response = &dto.ServerListResponse{
		Data: serverDTOs,
		Pagination: dto.PaginationDTO{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	// Cache the response
	if s.cacheRepo != nil {
		cached, err := json.Marshal(response)
		if err == nil {
			s.cacheRepo.Set(ctx, cacheKey, string(cached), s.cacheTTL)
		}
	}

	return response, nil
}

// Get server by its ID
func (s *ServerServiceImpl) GetServerByID(ctx context.Context, id int) (*dto.ServerDetailResponse, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf("server:%d", id)

	// Try to get from cache first
	var response *dto.ServerDetailResponse
	if s.cacheRepo != nil {
		cached, err := s.cacheRepo.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			err = json.Unmarshal([]byte(cached), &response)
			if err == nil {
				return response, nil
			}
		}
	}

	// Get from database
	server, err := s.serverRepo.GetServerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	if server == nil {
		return nil, fmt.Errorf("server not found")
	}

	// Convert to DTO
	serverDTO := s.convertModelToDTO(*server)

	response = &dto.ServerDetailResponse{
		Data: serverDTO,
	}

	// Cache the response
	if s.cacheRepo != nil {
		cached, err := json.Marshal(response)
		if err == nil {
			s.cacheRepo.Set(ctx, cacheKey, string(cached), s.cacheTTL)
		}
	}

	return response, nil
}

// Get all unique locations
func (s *ServerServiceImpl) GetLocations(ctx context.Context) ([]string, error) {
	// Generate cache key
	cacheKey := "locations"

	// Try to get from cache first
	if s.cacheRepo != nil {
		cached, err := s.cacheRepo.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var locations []string
			err = json.Unmarshal([]byte(cached), &locations)
			if err == nil {
				return locations, nil
			}
		}
	}

	// Get from database
	locations, err := s.serverRepo.GetLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	// Cache the response
	if s.cacheRepo != nil {
		cached, err := json.Marshal(locations)
		if err == nil {
			s.cacheRepo.Set(ctx, cacheKey, string(cached), s.cacheTTL)
		}
	}

	return locations, nil
}

// Get metrics about the servers
func (s *ServerServiceImpl) GetMetrics(ctx context.Context) (*dto.MetricsResponse, error) {
	// Generate cache key
	cacheKey := "metrics"

	// Try to get from cache first
	var response *dto.MetricsResponse
	if s.cacheRepo != nil {
		cached, err := s.cacheRepo.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			err = json.Unmarshal([]byte(cached), &response)
			if err == nil {
				return response, nil
			}
		}
	}

	// Get from database
	metrics, err := s.serverRepo.GetMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	// Convert to DTO
	response = &dto.MetricsResponse{
		TotalServers:   metrics.TotalServers,
		AveragePrice:   metrics.AveragePrice,
		LocationsCount: metrics.LocationsCount,
		LastUpdated:    metrics.LastUpdated,
	}

	// Cache the response
	if s.cacheRepo != nil {
		cached, err := json.Marshal(response)
		if err == nil {
			s.cacheRepo.Set(ctx, cacheKey, string(cached), s.cacheTTL)
		}
	}

	return response, nil
}

// Convert a DTO request to model filters
func (s *ServerServiceImpl) convertRequestToFilters(req dto.ServerListRequest) models.ServerFilters {
	return models.ServerFilters{
		Query:      req.Query,
		Location:   req.Location,
		RAMMin:     req.RAMMin,
		RAMMax:     req.RAMMax,
		RAMValues:  req.RAMValues,
		StorageMin: req.StorageMin,
		StorageMax: req.StorageMax,
		HDD:        req.HDD,
		PriceMin:   req.PriceMin,
		PriceMax:   req.PriceMax,
		Sort:       req.Sort,
		Page:       req.Page,
		PerPage:    req.PerPage,
	}
}

// Converts a model to DTO
func (s *ServerServiceImpl) convertModelToDTO(server models.Server) dto.ServerDTO {
	return dto.ServerDTO{
		ID:           server.ID,
		Model:        server.Model,
		CPU:          server.CPU,
		RAMGB:        server.RAMGB,
		HDD:          server.HDD,
		StorageGB:    server.StorageGB,
		LocationCity: server.LocationCity,
		LocationCode: server.LocationCode,
		PriceEUR:     server.PriceEUR,
		RawPrice:     server.RawPrice,
		RawRAM:       server.RawRAM,
		RawHDD:       server.RawHDD,
		CreatedAt:    server.CreatedAt,
	}
}

// Generate a cache key for the given filters
func (s *ServerServiceImpl) generateCacheKey(prefix string, filters models.ServerFilters) string {
	// Create a hash of the filters for the cache key
	filterData := map[string]interface{}{
		"query":       filters.Query,
		"location":    filters.Location,
		"ram_min":     filters.RAMMin,
		"ram_max":     filters.RAMMax,
		"storage_min": filters.StorageMin,
		"storage_max": filters.StorageMax,
		"hdd":         filters.HDD,
		"price_min":   filters.PriceMin,
		"price_max":   filters.PriceMax,
		"sort":        filters.Sort,
		"page":        filters.Page,
		"per_page":    filters.PerPage,
	}

	data, _ := json.Marshal(filterData)
	hash := md5.Sum(data)

	return fmt.Sprintf("%s:%x", prefix, hash)
}
