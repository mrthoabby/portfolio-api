package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/mrthoabby/portfolio-api/internal/common"
)

// RateLimiter implements a simple in-memory rate limiter per IP
type RateLimiter struct {
	requests  map[string][]time.Time
	syncMutex sync.RWMutex
	limit     int           // max requests
	window    time.Duration // time window
}

// NewRateLimiter creates a new rate limiter
// limit: maximum number of requests allowed
// window: time window for the limit
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rateLimiter := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup old entries every minute
	go rateLimiter.cleanup()

	return rateLimiter
}

func (instance *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		instance.syncMutex.Lock()
		now := time.Now()
		for ip, times := range instance.requests {
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) <= instance.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(instance.requests, ip)
			} else {
				instance.requests[ip] = valid
			}
		}
		instance.syncMutex.Unlock()
	}
}

func (instance *RateLimiter) isAllowed(ip string) bool {
	instance.syncMutex.Lock()
	defer instance.syncMutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-instance.window)

	// Filter requests within the window
	var validRequests []time.Time
	for _, t := range instance.requests[ip] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= instance.limit {
		instance.requests[ip] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	instance.requests[ip] = validRequests
	return true
}

// Limit returns a middleware that rate limits requests
func (instance *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if !instance.isAllowed(ip) {
			w.Header().Set("Retry-After", "60")
			common.RespondError(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests. Please try again later.", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
