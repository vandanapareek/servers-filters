package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Server represents a server record in the database
type Server struct {
	ID           int       `db:"id" json:"id"`
	Model        string    `db:"model" json:"model"`
	CPU          *string   `db:"cpu" json:"cpu"`
	RAMGB        *int      `db:"ram_gb" json:"ram_gb"`
	HDD          string    `db:"hdd" json:"hdd"`
	StorageGB    *int      `db:"storage_gb" json:"storage_gb"`
	LocationCity *string   `db:"location_city" json:"location_city"`
	LocationCode *string   `db:"location_code" json:"location_code"`
	PriceEUR     *float64  `db:"price_eur" json:"price_eur"`
	RawPrice     string    `db:"raw_price" json:"raw_price"`
	RawRAM       string    `db:"raw_ram" json:"raw_ram"`
	RawHDD       string    `db:"raw_hdd" json:"raw_hdd"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// ServerFilters represents the filter parameters for server queries
type ServerFilters struct {
	Query      string   `json:"query"`
	Location   []string `json:"location"`
	RAMMin     *int     `json:"ram_min"`
	RAMMax     *int     `json:"ram_max"`
	RAMValues  []int    `json:"ram_values"`
	StorageMin *int     `json:"storage_min"`
	StorageMax *int     `json:"storage_max"`
	HDD        string   `json:"hdd"`
	PriceMin   *float64 `json:"price_min"`
	PriceMax   *float64 `json:"price_max"`
	Sort       string   `json:"sort"`
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ServerListResponse represents the response for server list endpoint
type ServerListResponse struct {
	Data       []Server   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// ServerDetailResponse represents the response for server detail endpoint
type ServerDetailResponse struct {
	Data Server `json:"data"`
}

// SortOption represents a sort option
type SortOption struct {
	Field string `json:"field"`
	Order string `json:"order"` // "asc" or "desc"
}

// ParseSort parses the sort string into field and order
func ParseSort(sortStr string) SortOption {
	if sortStr == "" {
		return SortOption{Field: "id", Order: "asc"}
	}

	// Expected format: "field.order" (e.g., "price.asc", "ram.desc")
	parts := []string{"id", "asc"} // default
	if len(sortStr) > 0 {
		// Split by dot
		if dotIndex := len(sortStr) - 4; dotIndex > 0 && sortStr[dotIndex:] == ".asc" {
			parts = []string{sortStr[:dotIndex], "asc"}
		} else if dotIndex := len(sortStr) - 5; dotIndex > 0 && sortStr[dotIndex:] == ".desc" {
			parts = []string{sortStr[:dotIndex], "desc"}
		} else {
			// If no order specified, default to asc
			parts = []string{sortStr, "asc"}
		}
	}

	// Validate field
	validFields := map[string]bool{
		"id":            true,
		"model":         true,
		"cpu":           true,
		"ram_gb":        true,
		"storage_gb":    true,
		"location_city": true,
		"price_eur":     true,
		"created_at":    true,
	}

	if !validFields[parts[0]] {
		parts[0] = "id"
	}

	// Validate order
	if parts[1] != "asc" && parts[1] != "desc" {
		parts[1] = "asc"
	}

	return SortOption{Field: parts[0], Order: parts[1]}
}

// Value implements the driver.Valuer interface for JSON encoding
func (sf ServerFilters) Value() (driver.Value, error) {
	return json.Marshal(sf)
}

// Scan implements the sql.Scanner interface for JSON decoding
func (sf *ServerFilters) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, sf)
}

// ServerMetrics represents metrics about servers
type ServerMetrics struct {
	TotalServers   int64     `json:"total_servers"`
	MinPrice       float64   `json:"min_price"`
	MaxPrice       float64   `json:"max_price"`
	LocationsCount int64     `json:"locations_count"`
	LastUpdated    time.Time `json:"last_updated"`
}
