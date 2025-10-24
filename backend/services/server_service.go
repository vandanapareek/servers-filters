package services

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"

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

	// Get from database (pagination is already applied at database level)
	servers, total, err := s.serverRepo.GetServers(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}

	// Debug logging
	fmt.Printf("DEBUG: Page=%d, PerPage=%d, Servers returned=%d, Total=%d\n",
		req.Page, req.PerPage, len(servers), total)

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
		MinPrice:       metrics.MinPrice,
		MaxPrice:       metrics.MaxPrice,
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
	// Convert TB values to GB for database filtering
	var storageMinGB, storageMaxGB *int

	if req.StorageMin != nil {
		// Convert TB to GB: multiply by 1024
		gb := int(*req.StorageMin * 1024)
		storageMinGB = &gb
		fmt.Printf("🔥 TB CONVERSION: %.2f TB -> %d GB\n", *req.StorageMin, gb)
	}

	if req.StorageMax != nil {
		// Convert TB to GB: multiply by 1024
		gb := int(*req.StorageMax * 1024)
		storageMaxGB = &gb
		fmt.Printf("🔥 TB CONVERSION: %.2f TB -> %d GB\n", *req.StorageMax, gb)
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
		PriceMin:   req.PriceMin,
		PriceMax:   req.PriceMax,
		Sort:       req.Sort,
		Page:       req.Page,
		PerPage:    req.PerPage,
	}
}

// Converts a model to DTO
func (s *ServerServiceImpl) convertModelToDTO(server models.Server) dto.ServerDTO {
	storageDisplay := s.formatStorageDisplay(server.StorageGB)
	hddType := s.extractHDDType(server.HDD)

	return dto.ServerDTO{
		ID:             server.ID,
		Model:          server.Model,
		CPU:            server.CPU,
		RAMGB:          server.RAMGB,
		HDD:            server.HDD,
		HDDType:        hddType,
		StorageGB:      server.StorageGB,
		StorageDisplay: storageDisplay,
		LocationCity:   server.LocationCity,
		LocationCode:   server.LocationCode,
		PriceEUR:       server.PriceEUR,
		RawPrice:       server.RawPrice,
		RawRAM:         server.RawRAM,
		RawHDD:         server.RawHDD,
		CreatedAt:      server.CreatedAt,
	}
}

// formatStorageDisplay formats storage in GB to TB when appropriate
func (s *ServerServiceImpl) formatStorageDisplay(storageGB *int) string {
	if storageGB == nil {
		return ""
	}

	gb := *storageGB
	if gb >= 1024 {
		// Convert to TB
		tb := float64(gb) / 1024.0
		if tb == float64(int(tb)) {
			// If it's a whole number, don't show decimal
			return fmt.Sprintf("%.0fTB", tb)
		} else {
			// Show one decimal place
			return fmt.Sprintf("%.1fTB", tb)
		}
	}

	return fmt.Sprintf("%dGB", gb)
}

// extractHDDType extracts the disk type from HDD string
func (s *ServerServiceImpl) extractHDDType(hdd string) string {
	if hdd == "" {
		return ""
	}

	// Extract disk type from HDD string (e.g., "4x480GBSSD" -> "SSD")
	// Look for common disk types at the end of the string
	diskTypes := []string{"SSD", "SATA2", "SATA3", "NVMe", "HDD"}
	for _, diskType := range diskTypes {
		if strings.HasSuffix(strings.ToUpper(hdd), diskType) {
			return diskType
		}
	}

	// If no specific type found, try to extract the last word after numbers and GB/TB
	// Remove numbers and GB/TB units, then get the last part
	cleaned := strings.ToUpper(hdd)
	// Remove common patterns like "2x500GB" or "4x480GB"
	cleaned = strings.ReplaceAll(cleaned, "GB", "")
	cleaned = strings.ReplaceAll(cleaned, "TB", "")

	// Split by numbers and get the last non-empty part
	parts := strings.FieldsFunc(cleaned, func(c rune) bool {
		return c >= '0' && c <= '9'
	})
	if len(parts) > 0 {
		lastPart := strings.TrimSpace(parts[len(parts)-1])
		if lastPart != "" {
			return lastPart
		}
	}

	return hdd
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
