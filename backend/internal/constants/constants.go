package constants

// HTTP Status Codes
const (
	StatusInternalServerError = 500
)

// HTTP Error Messages
const (
	ErrorInternalServerError  = "Internal Server Error"
	ErrorFailedToGetServers   = "Failed to retrieve servers"
	ErrorFailedToGetLocations = "Failed to retrieve locations"
	ErrorFailedToGetMetrics   = "Failed to retrieve metrics"
)

// Pagination Constants
const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Storage Conversion Constants
const (
	// TB to GB conversion factor
	TBToGBMultiplier = 1024
)

// Cache Constants
const (
	DefaultCacheTTL   = 300 // 5 minutes in seconds
	MetricsCacheKey   = "metrics"
	LocationsCacheKey = "locations"
)

// Server Constants
const (
	DefaultShutdownTimeout = 30 // seconds
)

// CORS Constants
const (
	DefaultCORSMaxAge = 300 // 5 minutes
)
