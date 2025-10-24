package dto

import (
	"time"
)

// Server object for API responses
type ServerDTO struct {
	ID             int       `json:"id"`
	Model          string    `json:"model"`
	CPU            *string   `json:"cpu,omitempty"`
	RAMGB          *int      `json:"ram_gb,omitempty"`
	HDD            string    `json:"hdd"`
	HDDType        string    `json:"hdd_type,omitempty"`
	StorageGB      *int      `json:"storage_gb,omitempty"`
	StorageDisplay string    `json:"storage_display,omitempty"`
	LocationCity   *string   `json:"location_city,omitempty"`
	LocationCode   *string   `json:"location_code,omitempty"`
	PriceEUR       *float64  `json:"price_eur,omitempty"`
	RawPrice       string    `json:"raw_price"`
	RawRAM         string    `json:"raw_ram"`
	RawHDD         string    `json:"raw_hdd"`
	CreatedAt      time.Time `json:"created_at"`
}

// Request parameters for server list endpoint
type ServerListRequest struct {
	Query      string   `json:"query" form:"q"`
	Location   []string `json:"location" form:"location"`
	RAMMin     *int     `json:"ram_min" form:"ram_min"`
	RAMMax     *int     `json:"ram_max" form:"ram_max"`
	RAMValues  []int    `json:"ram_values" form:"ram_values"`
	StorageMin *float64 `json:"storage_min" form:"storage_min"`
	StorageMax *float64 `json:"storage_max" form:"storage_max"`
	HDD        string   `json:"hdd" form:"hdd"`
	PriceMin   *float64 `json:"price_min" form:"price_min"`
	PriceMax   *float64 `json:"price_max" form:"price_max"`
	Sort       string   `json:"sort" form:"sort"`
	Page       int      `json:"page" form:"page"`
	PerPage    int      `json:"per_page" form:"per_page"`
}

// Pagination metadata for API responses
type PaginationDTO struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Response for server list endpoint
type ServerListResponse struct {
	Data       []ServerDTO   `json:"data"`
	Pagination PaginationDTO `json:"pagination"`
}

// Response for server detail endpoint
type ServerDetailResponse struct {
	Data ServerDTO `json:"data"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// Health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// Metrics response
type MetricsResponse struct {
	TotalServers   int64     `json:"total_servers"`
	MinPrice       float64   `json:"min_price"`
	MaxPrice       float64   `json:"max_price"`
	LocationsCount int64     `json:"locations_count"`
	LastUpdated    time.Time `json:"last_updated"`
}
