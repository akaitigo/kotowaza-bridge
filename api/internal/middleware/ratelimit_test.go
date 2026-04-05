package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestIPRateLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:            rate.Every(100 * time.Millisecond), // 10/sec
		Burst:           3,
		CleanupInterval: 1 * time.Minute,
		MaxAge:          5 * time.Minute,
	}
	rl := NewIPRateLimiter(cfg)
	defer rl.Close()

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 3 requests should succeed (burst).
	for i := range 3 {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "request %d should succeed", i+1)
	}
}

func TestIPRateLimiter_RejectsExcessRequests(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:            rate.Every(1 * time.Hour), // very slow refill
		Burst:           2,
		CleanupInterval: 1 * time.Minute,
		MaxAge:          5 * time.Minute,
	}
	rl := NewIPRateLimiter(cfg)
	defer rl.Close()

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// First 2 requests exhaust the burst.
	for range 2 {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 3rd request should be rate limited.
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "6", w.Header().Get("Retry-After"))
	assert.Contains(t, w.Body.String(), "rate limit exceeded")
}

func TestIPRateLimiter_DifferentIPsHaveSeparateLimits(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:            rate.Every(1 * time.Hour),
		Burst:           1,
		CleanupInterval: 1 * time.Minute,
		MaxAge:          5 * time.Minute,
	}
	rl := NewIPRateLimiter(cfg)
	defer rl.Close()

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// IP 1: exhaust its limit.
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/chat", nil)
	req1.RemoteAddr = "10.0.0.1:1111"
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	w1b := httptest.NewRecorder()
	handler.ServeHTTP(w1b, req1)
	assert.Equal(t, http.StatusTooManyRequests, w1b.Code)

	// IP 2: should still be allowed.
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/chat", nil)
	req2.RemoteAddr = "10.0.0.2:2222"
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestExtractIP(t *testing.T) {
	t.Run("uses X-Real-IP when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", "1.2.3.4")
		req.RemoteAddr = "5.6.7.8:9999"
		assert.Equal(t, "1.2.3.4", extractIP(req))
	})

	t.Run("uses first X-Forwarded-For when X-Real-IP absent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
		req.RemoteAddr = "9.9.9.9:9999"
		assert.Equal(t, "1.2.3.4", extractIP(req))
	})

	t.Run("uses single X-Forwarded-For", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		req.RemoteAddr = "9.9.9.9:9999"
		assert.Equal(t, "1.2.3.4", extractIP(req))
	})

	t.Run("falls back to RemoteAddr", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		assert.Equal(t, "192.168.1.1", extractIP(req))
	})

	t.Run("handles RemoteAddr without port", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1"
		assert.Equal(t, "192.168.1.1", extractIP(req))
	})
}

func TestIPRateLimiter_Cleanup(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:            rate.Every(time.Second),
		Burst:           1,
		CleanupInterval: 50 * time.Millisecond,
		MaxAge:          50 * time.Millisecond,
	}
	rl := NewIPRateLimiter(cfg)
	defer rl.Close()

	// Add an entry.
	rl.getLimiter("10.0.0.1")

	rl.mu.Lock()
	assert.Equal(t, 1, len(rl.limiters))
	rl.mu.Unlock()

	// Wait for cleanup to run.
	time.Sleep(200 * time.Millisecond)

	rl.mu.Lock()
	assert.Equal(t, 0, len(rl.limiters), "stale entries should be cleaned up")
	rl.mu.Unlock()
}
