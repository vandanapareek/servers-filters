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

		filteredServers = append(filteredServers, server)
	}

	return filteredServers, int64(len(filteredServers)), nil
}

func (m *MockServerRepository) GetServerByID(ctx context.Context, id int) (*models.Server, error) {
	for _, server := range m.servers {
		if server.ID == id {
			return &server, nil
		}
	}
	return nil, nil
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

func TestServerService_GetServerByID(t *testing.T) {
	mockServer := models.Server{
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
	}

	mockRepo := &MockServerRepository{
		servers: []models.Server{mockServer},
	}
	mockCache := &MockCacheRepository{}

	service := NewServerService(mockRepo, mockCache, 300)

	// Test existing server
	response, err := service.GetServerByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetServerByID() error = %v", err)
	}

	if response.Data.ID != 1 {
		t.Errorf("GetServerByID() got ID = %d, want 1", response.Data.ID)
	}

	// Test non-existing server
	_, err = service.GetServerByID(context.Background(), 999)
	if err == nil {
		t.Error("GetServerByID() expected error for non-existing server")
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
		AveragePrice:   105.9,
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

	if metrics.AveragePrice != mockMetrics.AveragePrice {
		t.Errorf("GetMetrics() got %f average price, want %f", metrics.AveragePrice, mockMetrics.AveragePrice)
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
