package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       float64   // tokens per second
	burst      int       // maximum tokens in the bucket
	tokens     float64   // current tokens available
	lastUpdate time.Time // last time tokens were updated
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified rate (requests per second) and burst capacity
// rate must be positive and burst must be greater than zero
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	if rate <= 0 {
		panic("rate must be positive")
	}
	if burst <= 0 {
		panic("burst must be greater than zero")
	}

	return &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastUpdate: time.Now(),
	}
}

// refillTokens updates the token count based on elapsed time
// Must be called with the mutex held
func (rl *RateLimiter) refillTokens(now time.Time) {
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.tokens += elapsed * rl.rate
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}
	rl.lastUpdate = now
}

// Allow checks if a request is allowed based on the rate limit
// Returns true if the request is allowed, false otherwise
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refillTokens(time.Now())

	// Check if we have at least 1 token
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}

// Wait blocks until a request is allowed
func (rl *RateLimiter) Wait() {
	for {
		rl.mu.Lock()
		rl.refillTokens(time.Now())

		// Check if we have at least 1 token
		if rl.tokens >= 1.0 {
			rl.tokens -= 1.0
			rl.mu.Unlock()
			return
		}

		// Calculate exact wait time needed for 1 token
		tokensNeeded := 1.0 - rl.tokens
		waitTime := time.Duration(tokensNeeded/rl.rate*1000) * time.Millisecond
		rl.mu.Unlock()

		time.Sleep(waitTime)
	}
}

// GetTokens returns the current number of available tokens (for testing/monitoring)
// Note: This is a read-only operation and does not update the internal state
func (rl *RateLimiter) GetTokens() float64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()

	// Calculate tokens without modifying state
	tokens := rl.tokens + elapsed*rl.rate
	if tokens > float64(rl.burst) {
		tokens = float64(rl.burst)
	}

	return tokens
}
