package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	allowed := []string{"https://app.example.com", "http://localhost:3000"}
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})
	handler := CORS(allowed)(next)

	t.Run("echoes an allowed origin and calls next", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza", nil)
		req.Header.Set("Origin", "https://app.example.com")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.True(t, nextCalled)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "https://app.example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Values("Vary"), "Origin")
	})

	t.Run("does not reflect a disallowed origin but still serves the request", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza", nil)
		req.Header.Set("Origin", "https://evil.example.com")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.True(t, nextCalled)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Values("Vary"), "Origin")
	})

	t.Run("sets no allow-origin when the Origin header is absent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("answers a preflight from an allowed origin with 204 and CORS headers", func(t *testing.T) {
		nextCalled = false
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/chat", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.False(t, nextCalled)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	})

	t.Run("answers a preflight from a disallowed origin with 204 and no allow-origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/chat", nil)
		req.Header.Set("Origin", "https://evil.example.com")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}
