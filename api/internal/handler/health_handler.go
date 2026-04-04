package handler

import "net/http"

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health handles GET /api/v1/health
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
