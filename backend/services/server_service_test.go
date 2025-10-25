package services

import (
	"context"
	"testing"
	"time"

	"servers-filters/dto"
	"servers-filters/models"
)

// MockServerRepository implements ServerRepository for testing
type MockServerRepository struct {
	servers   []models.Server
	locations []string
	metrics   *models.ServerMetrics
}

func (m *MockServerRepository) GetServers(ctx context.Context, filters models.ServerFilters) ([]models.Server, int64, error) {
	// Apply basic filtering for testing
	filteredServers := make([]models.Server, 0)

	for _, server := range m.servers {
		// Apply RAM filter
		if filters.RAMMin != nil && server.RAMGB != nil && *server.RAMGB < *filters.RAMMin {
			continue
		}
		if filters.RAMMax != nil && server.RAMGB != nil && *server.RAMGB > *filters.RAMMax {
			continue
		}

		// Apply RAM values filter (exact match)
		if len(filters.RAMValues) > 0 && server.RAMGB != nil {
			found := false
			for _, ramValue := range filters.RAMValues {
				if *server.RAMGB == ramValue {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply storage filter (TB to GB conversion)
		if filters.StorageMin != nil && server.StorageGB != nil && *server.StorageGB < *filters.StorageMin {
			continue
		}
		if filters.StorageMax != nil && server.StorageGB != nil && *server.StorageGB > *filters.StorageMax {
			continue
		}

		filteredServers = append(filteredServers, server)
	}

	// Apply pagination
	total := int64(len(filteredServers))
	start := (filters.Page - 1) * filters.PerPage
	end := start + filters.PerPage

	if start >= len(filteredServers) {
		return []models.Server{}, total, nil
	}

	if end > len(filteredServers) {
		end = len(filteredServers)
	}

	return filteredServers[start:end], total, nil
}


func (m *MockServerRepository) GetServerCount(ctx context.Context, filters models.ServerFilters) (int64, error) {
	// Apply same filtering logic as GetServers
	filteredServers := make([]models.Server, 0)

	for _, server := range m.servers {
		// Apply RAM filter
		if filters.RAMMin != nil && server.RAMGB != nil && *server.RAMGB < *filters.RAMMin {
			continue
		}
		if filters.RAMMax != nil && server.RAMGB != nil && *server.RAMGB > *filters.RAMMax {
			continue
		}

		// Apply RAM values filter (exact match)
		if len(filters.RAMValues) > 0 && server.RAMGB != nil {
			found := false
			for _, ramValue := range filters.RAMValues {
				if *server.RAMGB == ramValue {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filteredServers = append(filteredServers, server)
	}

	return int64(len(filteredServers)), nil
}

func (m *MockServerRepository) GetLocations(ctx context.Context) ([]string, error) {
	return m.locations, nil
}

func (m *MockServerRepository) GetMetrics(ctx context.Context) (*models.ServerMetrics, error) {
	return m.metrics, nil
}

// MockCacheRepository implements CacheRepository for testing
type MockCacheRepository struct{}

func (m *MockCacheRepository) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value string, ttl int) error {
	return nil
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *MockCacheRepository) Clear(ctx context.Context) error {
	return nil
}

func TestServerService_GetServers(t *testing.T) {
	// Setup mock data
	mockServers := []models.Server{
		{
			ID:           1,
			Model:        "Dell R740",
			CPU:          stringPtr("Intel Xeon"),
			RAMGB:        intPtr(32),
			HDD:          "2x2TB SATA2",
			StorageGB:    intPtr(4096),
			LocationCity: stringPtr("Amsterdam"),
			LocationCode: stringPtr("AMS-01"),
			PriceEUR:     float64Ptr(89.0),
			RawPrice:     "€89.00",
			RawRAM:       "32GB DDR4",
			RawHDD:       "2x2TB SATA2",
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			Model:        "HP DL380",
			CPU:          stringPtr("AMD Ryzen"),
			RAMGB:        intPtr(16),
			HDD:          "1x500GB SSD",
			StorageGB:    intPtr(500),
			LocationCity: stringPtr("New York"),
			LocationCode: stringPtr("NY-01"),
			PriceEUR:     float64Ptr(75.5),
			RawPrice:     "€75.50",
			RawRAM:       "16GB DDR3",
			RawHDD:       "1x500GB SSD",
			CreatedAt:    time.Now(),
		},
	}

	mockRepo := &MockServerRepository{
		servers: mockServers,
	}
	mockCache := &MockCacheRepository{}

	service := NewServerService(mockRepo, mockCache, 300)

	// Test cases
	tests := []struct {
		name     string
		request  dto.ServerListRequest
		expected int
	}{
		{
			name: "Get all servers",
			request: dto.ServerListRequest{
				Page:    1,
				PerPage: 20,
			},
			expected: 2,
		},
		{
			name: "Get servers with pagination",
			request: dto.ServerListRequest{
				Page:    1,
				PerPage: 1,
			},
			expected: 1,
		},
		{
			name: "Get servers with RAM filter",
			request: dto.ServerListRequest{
				RAMMin:  intPtr(32),
				Page:    1,
				PerPage: 20,
			},
			expected: 1,
		},
		{
			name: "Get servers with storage filter (TB values)",
			request: dto.ServerListRequest{
				StorageMin: float64Ptr(1.0), // 1TB
				StorageMax: float64Ptr(5.0), // 5TB
				Page:       1,
				PerPage:    20,
			},
			expected: 1, // Only the 4TB server should match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.GetServers(context.Background(), tt.request)
			if err != nil {
				t.Fatalf("GetServers() error = %v", err)
			}

			if len(response.Data) != tt.expected {
				t.Errorf("GetServers() got %d servers, want %d", len(response.Data), tt.expected)
			}

			if response.Pagination.Page != tt.request.Page {
				t.Errorf("GetServers() pagination page = %d, want %d", response.Pagination.Page, tt.request.Page)
			}
		})
	}
}


func TestServerService_GetLocations(t *testing.T) {
	mockLocations := []string{"Amsterdam", "New York", "London"}

	mockRepo := &MockServerRepository{
		locations: mockLocations,
	}
	mockCache := &MockCacheRepository{}

	service := NewServerService(mockRepo, mockCache, 300)

	locations, err := service.GetLocations(context.Background())
	if err != nil {
		t.Fatalf("GetLocations() error = %v", err)
	}

	if len(locations) != len(mockLocations) {
		t.Errorf("GetLocations() got %d locations, want %d", len(locations), len(mockLocations))
	}
}

func TestServerService_GetMetrics(t *testing.T) {
	mockMetrics := &models.ServerMetrics{
		TotalServers:   5,
		MinPrice:       35.99,
		MaxPrice:       4662.99,
		LocationsCount: 4,
		LastUpdated:    time.Now(),
	}

	mockRepo := &MockServerRepository{
		metrics: mockMetrics,
	}
	mockCache := &MockCacheRepository{}

	service := NewServerService(mockRepo, mockCache, 300)

	metrics, err := service.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("GetMetrics() error = %v", err)
	}

	if metrics.TotalServers != mockMetrics.TotalServers {
		t.Errorf("GetMetrics() got %d total servers, want %d", metrics.TotalServers, mockMetrics.TotalServers)
	}

	if metrics.MinPrice != mockMetrics.MinPrice {
		t.Errorf("GetMetrics() got %f min price, want %f", metrics.MinPrice, mockMetrics.MinPrice)
	}

	if metrics.MaxPrice != mockMetrics.MaxPrice {
		t.Errorf("GetMetrics() got %f max price, want %f", metrics.MaxPrice, mockMetrics.MaxPrice)
	}
}

func TestServerService_ExtractHDDType(t *testing.T) {
	service := &ServerServiceImpl{}

	tests := []struct {
		name     string
		hdd      string
		expected string
	}{
		{
			name:     "SSD extraction",
			hdd:      "4x480GBSSD",
			expected: "SSD",
		},
		{
			name:     "SATA2 extraction",
			hdd:      "2x2TBSATA2",
			expected: "SATA2",
		},
		{
			name:     "SATA3 extraction",
			hdd:      "1x1TBSATA3",
			expected: "SATA3",
		},
		{
			name:     "NVMe extraction",
			hdd:      "2x500GBNVMe",
			expected: "NVME",
		},
		{
			name:     "HDD extraction",
			hdd:      "1x1TBHDD",
			expected: "HDD",
		},
		{
			name:     "Case insensitive",
			hdd:      "4x480GBssd",
			expected: "SSD",
		},
		{
			name:     "Empty string",
			hdd:      "",
			expected: "",
		},
		{
			name:     "Unknown type",
			hdd:      "2x500GBUnknown",
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractHDDType(tt.hdd)
			if result != tt.expected {
				t.Errorf("extractHDDType(%s) = %s, want %s", tt.hdd, result, tt.expected)
			}
		})
	}
}

func TestServerService_FormatStorageDisplay(t *testing.T) {
	service := &ServerServiceImpl{}

	tests := []struct {
		name     string
		storage  *int
		expected string
	}{
		{
			name:     "GB storage",
			storage:  intPtr(500),
			expected: "500GB",
		},
		{
			name:     "TB storage whole number",
			storage:  intPtr(4096),
			expected: "4TB",
		},
		{
			name:     "TB storage with decimal",
			storage:  intPtr(1920),
			expected: "1.9TB",
		},
		{
			name:     "Nil storage",
			storage:  nil,
			expected: "",
		},
		{
			name:     "Exact 1TB",
			storage:  intPtr(1024),
			expected: "1TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.formatStorageDisplay(tt.storage)
			if result != tt.expected {
				t.Errorf("formatStorageDisplay(%v) = %s, want %s", tt.storage, result, tt.expected)
			}
		})
	}
}

func TestServerService_ConvertModelToDTO(t *testing.T) {
	service := &ServerServiceImpl{}

	server := models.Server{
		ID:           1,
		Model:        "Dell R740",
		CPU:          stringPtr("Intel Xeon"),
		RAMGB:        intPtr(32),
		HDD:          "2x2TBSATA2",
		StorageGB:    intPtr(4096),
		LocationCity: stringPtr("Amsterdam"),
		LocationCode: stringPtr("AMS-01"),
		PriceEUR:     float64Ptr(89.0),
		RawPrice:     "€89.00",
		RawRAM:       "32GB DDR4",
		RawHDD:       "2x2TB SATA2",
		CreatedAt:    time.Now(),
	}

	dto := service.convertModelToDTO(server)

	// Test basic fields
	if dto.ID != server.ID {
		t.Errorf("convertModelToDTO() ID = %d, want %d", dto.ID, server.ID)
	}

	if dto.Model != server.Model {
		t.Errorf("convertModelToDTO() Model = %s, want %s", dto.Model, server.Model)
	}

	// Test HDD type extraction
	if dto.HDDType != "SATA2" {
		t.Errorf("convertModelToDTO() HDDType = %s, want SATA2", dto.HDDType)
	}

	// Test storage display formatting
	if dto.StorageDisplay != "4TB" {
		t.Errorf("convertModelToDTO() StorageDisplay = %s, want 4TB", dto.StorageDisplay)
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
