package handler

import (
	"context"
	"log"
	"net/http"
	"time"
)

// readinessTimeout bounds the database ping performed by the readiness check so
// a hung database cannot block the health endpoint indefinitely.
const readinessTimeout = 2 * time.Second

// Pinger is the subset of the database pool required for readiness checks.
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthHandler handles liveness and readiness checks.
type HealthHandler struct {
	db Pinger
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Live handles GET /api/v1/health and reports process liveness without
// touching downstream dependencies. It always returns 200 while the process is
// running so orchestrators do not restart a healthy process during a transient
// database outage.
func (h *HealthHandler) Live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

// Ready handles GET /api/v1/health/ready and reports whether the service can
// serve traffic by verifying database connectivity. It returns 503 when the
// database is unreachable so load balancers stop routing requests to this
// instance.
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		// Log the underlying cause server-side but return a generic message so
		// the public endpoint does not leak connection details.
		log.Printf("readiness check failed: database ping: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, HealthResponse{Status: "unavailable", Error: "database unreachable"})
		return
	}

	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
