package services

import (
	"context"
	"fmt"

	"servers-filters/dto"
	"servers-filters/internal/constants"
	"servers-filters/models"
	"servers-filters/repository"
)

// Implement ServerService
type ServerServiceImpl struct {
	serverRepo repository.ServerRepository
}

// Create new server service
func NewServerService(serverRepo repository.ServerRepository) ServerService {
	return &ServerServiceImpl{
		serverRepo: serverRepo,
	}
}

// Get servers with filters and pagination
func (s *ServerServiceImpl) GetServers(ctx context.Context, req dto.ServerListRequest) (*dto.ServerListResponse, error) {
	if req.Page <= 0 {
		req.Page = constants.DefaultPage
	}
	if req.PerPage <= 0 {
		req.PerPage = constants.DefaultPerPage
	}
	if req.PerPage > constants.MaxPerPage {
		req.PerPage = constants.MaxPerPage
	}

	// convert request to model filters
	filters := s.convertRequestToFilters(req)

	// get from database
	servers, total, err := s.serverRepo.GetServers(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}

	// convert to DTOs
	serverDTOs := make([]dto.ServerDTO, len(servers))
	for i, server := range servers {
		serverDTOs[i] = s.convertModelToDTO(server)
	}

	// calculate pagination
	totalPages := int((total + int64(req.PerPage) - 1) / int64(req.PerPage))

	response := &dto.ServerListResponse{
		Data: serverDTOs,
		Pagination: dto.PaginationDTO{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

// Get all unique locations
func (s *ServerServiceImpl) GetLocations(ctx context.Context) ([]string, error) {
	// Get from database
	locations, err := s.serverRepo.GetLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	return locations, nil
}

// Get metrics about the servers
func (s *ServerServiceImpl) GetMetrics(ctx context.Context) (*dto.MetricsResponse, error) {
	// get from database
	metrics, err := s.serverRepo.GetMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	// convert to DTO
	response := &dto.MetricsResponse{
		TotalServers:   metrics.TotalServers,
		MinPrice:       metrics.MinPrice,
		MaxPrice:       metrics.MaxPrice,
		LocationsCount: metrics.LocationsCount,
		LastUpdated:    metrics.LastUpdated,
	}

	return response, nil
}

// Convert DTO request to model filters
func (s *ServerServiceImpl) convertRequestToFilters(req dto.ServerListRequest) models.ServerFilters {
	// Convert TB to GB for database filtering
	var storageMinGB, storageMaxGB *int

	if req.StorageMin != nil {
		gb := int(*req.StorageMin * constants.TBToGBMultiplier)
		storageMinGB = &gb
	}

	if req.StorageMax != nil {
		gb := int(*req.StorageMax * constants.TBToGBMultiplier)
		storageMaxGB = &gb
	}

	return models.ServerFilters{
		Query:      req.Query,
		Location:   req.Location,
		RAMMin:     req.RAMMin,
		RAMMax:     req.RAMMax,
		RAMValues:  req.RAMValues,
		StorageMin: storageMinGB,
		StorageMax: storageMaxGB,
		HDD:        req.HDD,
		Sort:       req.Sort,
		Page:       req.Page,
		PerPage:    req.PerPage,
	}
}

// Convert model to DTO
func (s *ServerServiceImpl) convertModelToDTO(server models.Server) dto.ServerDTO {
	storageDisplay := s.formatStorageDisplay(server.HDDGB)

	var hddType string
	if server.HDDType != nil {
		hddType = *server.HDDType
	}

	return dto.ServerDTO{
		ID:             server.ID,
		Model:          server.Model,
		CPU:            server.CPU,
		RAMGB:          server.RAMGB,
		HDDGB:          server.HDDGB,
		HDDType:        hddType,
		StorageDisplay: storageDisplay,
		Location:       server.Location,
		LocationCode:   server.LocationCode,
		Price:          server.Price,
		RawPrice:       server.RawPrice,
		RawHDD:         server.RawHDD,
		RawRAM:         server.RawRAM,
		CreatedAt:      server.CreatedAt,
		UpdatedAt:      server.UpdatedAt,
	}
}

// format storage in GB to TB
func (s *ServerServiceImpl) formatStorageDisplay(storageGB *int) string {
	if storageGB == nil {
		return ""
	}

	gb := *storageGB
	if gb >= 1000 {
		// convert to TB
		tb := float64(gb) / 1000.0
		if tb == float64(int(tb)) {
			// if it's a whole number, don't show decimal
			return fmt.Sprintf("%.0fTB", tb)
		} else {
			// show one decimal place
			return fmt.Sprintf("%.1fTB", tb)
		}
	}

	return fmt.Sprintf("%dGB", gb)
}
