package services

import (
	"context"
	"testing"
	"time"

	"servers-filters/dto"
	"servers-filters/models"
)

// implement ServerRepository for testing
type MockServerRepository struct {
	servers   []models.Server
	locations []string
	metrics   *models.ServerMetrics
}

func (m *MockServerRepository) GetServers(ctx context.Context, filters models.ServerFilters) ([]models.Server, int64, error) {
	// Apply basic filtering
	filteredServers := make([]models.Server, 0)

	for _, server := range m.servers {
		// Apply RAM filter
		if filters.RAMMin != nil && server.RAMGB != nil && *server.RAMGB < *filters.RAMMin {
			continue
		}
		if filters.RAMMax != nil && server.RAMGB != nil && *server.RAMGB > *filters.RAMMax {
			continue
		}

		// Apply RAM values filter
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

		// Apply storage filter
		if filters.StorageMin != nil && server.HDDGB != nil && *server.HDDGB < *filters.StorageMin {
			continue
		}
		if filters.StorageMax != nil && server.HDDGB != nil && *server.HDDGB > *filters.StorageMax {
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
	filteredServers := make([]models.Server, 0)

	for _, server := range m.servers {
		if filters.RAMMin != nil && server.RAMGB != nil && *server.RAMGB < *filters.RAMMin {
			continue
		}
		if filters.RAMMax != nil && server.RAMGB != nil && *server.RAMGB > *filters.RAMMax {
			continue
		}

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

func TestServerService_GetServers(t *testing.T) {
	// Setup mock data
	mockServers := []models.Server{
		{
			ID:           1,
			Model:        "Dell R740",
			CPU:          stringPtr("Intel Xeon"),
			RAMGB:        intPtr(32),
			HDDGB:        intPtr(4000),
			HDDType:      stringPtr("SATA2"),
			Location:     stringPtr("Amsterdam"),
			LocationCode: stringPtr("AMS-01"),
			Price:        float64Ptr(89.0),
			RawPrice:     "€89.00",
			RawRAM:       "32GB DDR4",
			RawHDD:       "2x2TB SATA2",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           2,
			Model:        "HP DL380",
			CPU:          stringPtr("AMD Ryzen"),
			RAMGB:        intPtr(16),
			HDDGB:        intPtr(500),
			HDDType:      stringPtr("SSD"),
			Location:     stringPtr("New York"),
			LocationCode: stringPtr("NY-01"),
			Price:        float64Ptr(75.5),
			RawPrice:     "€75.50",
			RawRAM:       "16GB DDR3",
			RawHDD:       "1x500GB SSD",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockRepo := &MockServerRepository{
		servers: mockServers,
	}

	service := NewServerService(mockRepo)

	// test cases
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

func TestServerService_GetLocations(t *testing.T) {
	mockLocations := []string{"Amsterdam", "New York", "London"}

	mockRepo := &MockServerRepository{
		locations: mockLocations,
	}

	service := NewServerService(mockRepo)

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

	service := NewServerService(mockRepo)

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
			storage:  intPtr(4000),
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
			storage:  intPtr(1000),
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
		HDDGB:        intPtr(4000),
		HDDType:      stringPtr("SATA2"),
		Location:     stringPtr("Amsterdam"),
		LocationCode: stringPtr("AMS-01"),
		Price:        float64Ptr(89.0),
		RawPrice:     "€89.00",
		RawRAM:       "32GB DDR4",
		RawHDD:       "2x2TB SATA2",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	dto := service.convertModelToDTO(server)

	// Test basic fields
	if dto.ID != server.ID {
		t.Errorf("convertModelToDTO() ID = %d, want %d", dto.ID, server.ID)
	}

	if dto.Model != server.Model {
		t.Errorf("convertModelToDTO() Model = %s, want %s", dto.Model, server.Model)
	}

	// Test HDD type from database
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
