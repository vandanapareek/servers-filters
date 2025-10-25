package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"servers-filters/dto"
	"servers-filters/internal/constants"
	"servers-filters/internal/logger"
	"servers-filters/services"

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
		StorageMin: parseFloatParam(r.URL.Query().Get("storage_min")),
		StorageMax: parseFloatParam(r.URL.Query().Get("storage_max")),
		HDD:        r.URL.Query().Get("hdd"),
		Sort:       r.URL.Query().Get("sort"),
		Page:       parseIntParamWithDefault(r.URL.Query().Get("page"), constants.DefaultPage),
		PerPage:    parseIntParamWithDefault(r.URL.Query().Get("per_page"), constants.DefaultPerPage),
	}

	// Debug: Check what we receive from API
	fmt.Println("🔥 API DEBUG: Starting handler")
	fmt.Printf("🔥 API DEBUG: storage_max raw='%s'\n", r.URL.Query().Get("storage_max"))
	fmt.Printf("🔥 API DEBUG: storage_max parsed=%v\n", req.StorageMax)

	// Get servers
	response, err := h.serverService.GetServers(r.Context(), req)
	if err != nil {
		logger.GetLogger().WithError(err).Error(constants.ErrorFailedToGetServers)
		render.Status(r, constants.StatusInternalServerError)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   constants.ErrorInternalServerError,
			Message: constants.ErrorFailedToGetServers,
			Code:    constants.StatusInternalServerError,
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
		logger.GetLogger().WithError(err).Error(constants.ErrorFailedToGetLocations)
		render.Status(r, constants.StatusInternalServerError)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   constants.ErrorInternalServerError,
			Message: constants.ErrorFailedToGetLocations,
			Code:    constants.StatusInternalServerError,
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
		logger.GetLogger().WithError(err).Error(constants.ErrorFailedToGetMetrics)
		render.Status(r, constants.StatusInternalServerError)
		render.JSON(w, r, dto.ErrorResponse{
			Error:   constants.ErrorInternalServerError,
			Message: constants.ErrorFailedToGetMetrics,
			Code:    constants.StatusInternalServerError,
		})
		return
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
