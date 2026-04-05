package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiterConfig holds configuration for the rate limiter.
type RateLimiterConfig struct {
	// Rate is the number of requests allowed per second.
	Rate rate.Limit
	// Burst is the maximum number of requests allowed in a burst.
	Burst int
	// CleanupInterval is how often stale entries are removed.
	CleanupInterval time.Duration
	// MaxAge is how long an IP entry is kept after last seen.
	MaxAge time.Duration
}

// DefaultChatRateLimiterConfig returns the default rate limiter config for chat endpoints.
// 10 requests per minute per IP, burst of 3.
func DefaultChatRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Rate:            rate.Every(6 * time.Second), // 10 per minute
		Burst:           3,
		CleanupInterval: 5 * time.Minute,
		MaxAge:          10 * time.Minute,
	}
}

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter tracks rate limiters per IP address.
type IPRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	config   RateLimiterConfig
	done     chan struct{}
}

// NewIPRateLimiter creates a new IP-based rate limiter and starts a background cleanup goroutine.
func NewIPRateLimiter(cfg RateLimiterConfig) *IPRateLimiter {
	rl := &IPRateLimiter{
		limiters: make(map[string]*ipLimiter),
		config:   cfg,
		done:     make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

// Close stops the background cleanup goroutine.
func (rl *IPRateLimiter) Close() {
	close(rl.done)
}

func (rl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.limiters[ip]
	if !exists {
		entry = &ipLimiter{
			limiter:  rate.NewLimiter(rl.config.Rate, rl.config.Burst),
			lastSeen: time.Now(),
		}
		rl.limiters[ip] = entry
	} else {
		entry.lastSeen = time.Now()
	}

	return entry.limiter
}

func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-rl.done:
			return
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, entry := range rl.limiters {
				if now.Sub(entry.lastSeen) > rl.config.MaxAge {
					delete(rl.limiters, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// Middleware returns an HTTP middleware that enforces the rate limit.
func (rl *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)

		limiter := rl.getLimiter(ip)
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "6")
			w.WriteHeader(http.StatusTooManyRequests)
			// Write JSON error response inline to avoid circular dependency with handler package.
			_, _ = w.Write([]byte(`{"error":"rate limit exceeded, please try again later"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// extractIP extracts the client IP from the request.
// It checks X-Forwarded-For and X-Real-IP headers before falling back to RemoteAddr.
func extractIP(r *http.Request) string {
	// Check X-Real-IP first (most specific single IP).
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Check X-Forwarded-For (first IP is the client).
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		for i := range len(xff) {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Fall back to RemoteAddr.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
