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
	HDDGB          *int      `json:"hdd_gb,omitempty"`
	HDDType        string    `json:"hdd_type,omitempty"`
	StorageDisplay string    `json:"storage_display,omitempty"`
	Location       *string   `json:"location,omitempty"`
	LocationCode   *string   `json:"location_code,omitempty"`
	Price          *float64  `json:"price,omitempty"`
	RawPrice       string    `json:"raw_price"`
	RawHDD         string    `json:"raw_hdd"`
	RawRAM         string    `json:"raw_ram"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
	Sort       string   `json:"sort" form:"sort"`
	Page       int      `json:"page" form:"page"`
	PerPage    int      `json:"per_page" form:"per_page"`
}

// Pagination Object for API responses
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

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// Metrics response
type MetricsResponse struct {
	TotalServers   int64     `json:"total_servers"`
	MinPrice       float64   `json:"min_price"`
	MaxPrice       float64   `json:"max_price"`
	LocationsCount int64     `json:"locations_count"`
	LastUpdated    time.Time `json:"last_updated"`
}
