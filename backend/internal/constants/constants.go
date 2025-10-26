package constants

const (
	StatusInternalServerError = 500
)

const (
	ErrorInternalServerError  = "Internal Server Error"
	ErrorFailedToGetServers   = "Failed to retrieve servers"
	ErrorFailedToGetLocations = "Failed to retrieve locations"
	ErrorFailedToGetMetrics   = "Failed to retrieve metrics"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

const (
	TBToGBMultiplier = 1024
)

const (
	DefaultShutdownTimeout = 30 // seconds
)

const (
	DefaultCORSMaxAge = 300 // 5 minutes
)
