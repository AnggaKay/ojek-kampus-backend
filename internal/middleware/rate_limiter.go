package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/labstack/echo/v4"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup old entries every minute
	go rl.cleanup()

	return rl
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, timestamps := range rl.requests {
			// Remove timestamps older than window
			filtered := []time.Time{}
			for _, ts := range timestamps {
				if now.Sub(ts) < rl.window {
					filtered = append(filtered, ts)
				}
			}
			if len(filtered) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = filtered
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get existing requests
	timestamps := rl.requests[key]

	// Filter out old timestamps
	filtered := []time.Time{}
	for _, ts := range timestamps {
		if now.Sub(ts) < rl.window {
			filtered = append(filtered, ts)
		}
	}

	// Check limit
	if len(filtered) >= rl.limit {
		return false
	}

	// Add new request
	filtered = append(filtered, now)
	rl.requests[key] = filtered

	return true
}

// RateLimit limits requests per IP
// Example: RateLimit(5, 15*time.Minute) = max 5 requests per 15 minutes
func RateLimit(limit int, window time.Duration) echo.MiddlewareFunc {
	limiter := newRateLimiter(limit, window)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			if !limiter.allow(ip) {
				return c.JSON(http.StatusTooManyRequests, dto.ErrorResponse(
					"RATE_LIMIT_EXCEEDED",
					"Too many requests. Please try again later.",
				))
			}

			return next(c)
		}
	}
}
