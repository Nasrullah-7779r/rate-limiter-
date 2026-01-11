package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter implements a sliding window log rate limiter
type RateLimiter struct {
	maxRequests int           // maximum requests allowed in the window
	window      time.Duration // time window duration
	requests    []time.Time   // log of request timestamps
	mu          sync.Mutex
}

// NewRateLimiter creates a new sliding window log rate limiter
// maxRequests: maximum number of requests allowed in the time window
// window: time window duration (e.g., 1 second, 1 minute)
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	if maxRequests <= 0 {
		panic("maxRequests must be positive")
	}
	if window <= 0 {
		panic("window must be positive")
	}

	return &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0, maxRequests),
	}
}

// cleanOldRequests removes requests that are outside the current time window
// Must be called with the mutex held
func (rl *RateLimiter) cleanOldRequests(now time.Time) {
	cutoff := now.Add(-rl.window)

	// Find the first request that is still within the window
	validIdx := 0
	for validIdx < len(rl.requests) && rl.requests[validIdx].Before(cutoff) {
		validIdx++
	}

	// Keep only the valid requests
	if validIdx > 0 {
		rl.requests = rl.requests[validIdx:]
	}
}

// Allow checks if a request is allowed based on the rate limit
// Returns true if the request is allowed, false otherwise
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.cleanOldRequests(now)

	// Check if we have capacity for another request
	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}

	return false
}

// Wait blocks until a request is allowed
func (rl *RateLimiter) Wait() {
	for {
		rl.mu.Lock()
		now := time.Now()
		rl.cleanOldRequests(now)

		// Check if we have capacity for another request
		if len(rl.requests) < rl.maxRequests {
			rl.requests = append(rl.requests, now)
			rl.mu.Unlock()
			return
		}

		// Calculate wait time: oldest request time + window - now
		waitUntil := rl.requests[0].Add(rl.window)
		waitTime := waitUntil.Sub(now)
		rl.mu.Unlock()

		if waitTime > 0 {
			time.Sleep(waitTime)
		}
	}
}

// GetRequestCount returns the current number of requests in the window (for testing/monitoring)
func (rl *RateLimiter) GetRequestCount() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.cleanOldRequests(time.Now())
	return len(rl.requests)
}
