package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"servers-filters/dto"
	"servers-filters/internal/logger"
	"servers-filters/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// ServerHandler handles server-related HTTP requests
type ServerHandler struct {
	serverService services.ServerService
}

// Creates a new server handler
func NewServerHandler(serverService services.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

// GetServers handles GET /servers
func (h *ServerHandler) GetServers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := dto.ServerListRequest{
		Query:      r.URL.Query().Get("q"),
		Location:   parseLocationParam(r.URL.Query().Get("location")),
		RAMMin:     parseIntParam(r.URL.Query().Get("ram_min")),
		RAMMax:     parseIntParam(r.URL.Query().Get("ram_max")),
		RAMValues:  parseIntArrayParam(r.URL.Query().Get("ram_values")),
		StorageMin: parseIntParam(r.URL.Query().Get("storage_min")),
		StorageMax: parseIntParam(r.URL.Query().Get("storage_max")),
		HDD:        r.URL.Query().Get("hdd"),
		PriceMin:   parseFloatParam(r.URL.Query().Get("price_min")),
		PriceMax:   parseFloatParam(r.URL.Query().Get("price_max")),
		Sort:       r.URL.Query().Get("sort"),
		Page:       parseIntParamWithDefault(r.URL.Query().Get("page"), 1),
		PerPage:    parseIntParamWithDefault(r.URL.Query().Get("per_page"), 20),
	}

	// Get servers
	response, err := h.serverService.GetServers(r.Context(), req)
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to get servers")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve servers",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	render.JSON(w, r, response)
}

// GetServerByID handles GET /servers/{id}
func (h *ServerHandler) GetServerByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid server ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Get server
	response, err := h.serverService.GetServerByID(r.Context(), id)
	if err != nil {
		logger.GetLogger().WithError(err).WithField("id", id).Error("Failed to get server")
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   "Not Found",
			Message: "Server not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	render.JSON(w, r, response)
}

// GetLocations handles GET /locations
func (h *ServerHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	// Get locations
	locations, err := h.serverService.GetLocations(r.Context())
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to get locations")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve locations",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"data": locations,
	})
}

// GetMetrics handles GET /metrics
func (h *ServerHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// Get metrics
	response, err := h.serverService.GetMetrics(r.Context())
	if err != nil {
		logger.GetLogger().WithError(err).Error("Failed to get metrics")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve metrics",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	render.JSON(w, r, response)
}

// Health handles GET /health
func (h *ServerHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	render.JSON(w, r, response)
}

// parseLocationParam parses comma-separated location values
func parseLocationParam(locationStr string) []string {
	if locationStr == "" {
		return nil
	}

	locations := strings.Split(locationStr, ",")
	var result []string
	for _, loc := range locations {
		loc = strings.TrimSpace(loc)
		if loc != "" {
			result = append(result, loc)
		}
	}

	return result
}

// parseIntParam parses an integer parameter
func parseIntParam(param string) *int {
	if param == "" {
		return nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil
	}

	return &val
}

// parseIntParamWithDefault parses an integer parameter with a default value
func parseIntParamWithDefault(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return val
}

// parseFloatParam parses a float parameter
func parseFloatParam(param string) *float64 {
	if param == "" {
		return nil
	}

	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return nil
	}

	return &val
}

// Parses a comma-separated integer array parameter
func parseIntArrayParam(param string) []int {
	if param == "" {
		return nil
	}

	parts := strings.Split(param, ",")
	var result []int

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			if val, err := strconv.Atoi(part); err == nil {
				result = append(result, val)
			}
		}
	}

	return result
}
