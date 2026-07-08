package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPinger struct {
	err error
}

func (s stubPinger) Ping(_ context.Context) error {
	return s.err
}

func TestHealthHandler_Live(t *testing.T) {
	h := NewHealthHandler(stubPinger{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	h.Live(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HealthResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "ok", resp.Status)
}

func TestHealthHandler_Ready(t *testing.T) {
	t.Run("returns 200 when the database is reachable", func(t *testing.T) {
		h := NewHealthHandler(stubPinger{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
		w := httptest.NewRecorder()

		h.Ready(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp HealthResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("returns 503 without leaking the ping error when the database is unreachable", func(t *testing.T) {
		h := NewHealthHandler(stubPinger{err: errors.New("dial tcp 10.0.0.1:5432: connection refused")})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
		w := httptest.NewRecorder()

		h.Ready(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var resp HealthResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, "unavailable", resp.Status)
		assert.NotContains(t, resp.Error, "connection refused")
		assert.NotContains(t, resp.Error, "10.0.0.1")
	})
}
