package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, time.Second)
	if rl.maxRequests != 10 {
		t.Errorf("Expected maxRequests 10, got %d", rl.maxRequests)
	}
	if rl.window != time.Second {
		t.Errorf("Expected window 1s, got %v", rl.window)
	}
	if len(rl.requests) != 0 {
		t.Errorf("Expected initial requests 0, got %d", len(rl.requests))
	}
}

func TestNewRateLimiter_InvalidMaxRequests(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero maxRequests")
		}
	}()
	NewRateLimiter(0, time.Second)
}

func TestNewRateLimiter_NegativeMaxRequests(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative maxRequests")
		}
	}()
	NewRateLimiter(-1, time.Second)
}

func TestNewRateLimiter_InvalidWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero window")
		}
	}()
	NewRateLimiter(10, 0)
}

func TestNewRateLimiter_NegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative window")
		}
	}()
	NewRateLimiter(10, -1*time.Second)
}

func TestRateLimiter_Allow_InitialRequests(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)

	// Should allow up to maxRequests immediately
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be denied (limit reached)
	if rl.Allow() {
		t.Error("Request should be denied after limit reached")
	}
}

func TestRateLimiter_Allow_WindowSliding(t *testing.T) {
	rl := NewRateLimiter(2, 200*time.Millisecond)

	// Use up the limit
	if !rl.Allow() {
		t.Error("First request should be allowed")
	}
	if !rl.Allow() {
		t.Error("Second request should be allowed")
	}

	// Should be denied immediately
	if rl.Allow() {
		t.Error("Third request should be denied")
	}

	// Wait for the window to slide (oldest request to expire)
	time.Sleep(250 * time.Millisecond)

	// Should be allowed now
	if !rl.Allow() {
		t.Error("Request should be allowed after window slides")
	}
}

func TestRateLimiter_Allow_RateLimit(t *testing.T) {
	maxRequests := 5
	window := time.Second
	rl := NewRateLimiter(maxRequests, window)

	// Count allowed requests over 2 seconds
	start := time.Now()
	allowed := 0
	for time.Since(start) < 2*time.Second {
		if rl.Allow() {
			allowed++
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Should allow approximately 10 requests over 2 seconds (5 per second)
	if allowed < 8 || allowed > 12 {
		t.Errorf("Expected ~10 requests allowed in 2 seconds, got %d", allowed)
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(50, time.Second)
	var wg sync.WaitGroup
	allowed := 0
	var mu sync.Mutex

	// Launch 100 concurrent goroutines
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should allow at most maxRequests
	if allowed != 50 {
		t.Errorf("Expected exactly 50 requests allowed, got %d", allowed)
	}
}

func TestRateLimiter_GetRequestCount(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	// Initial count should be 0
	count := rl.GetRequestCount()
	if count != 0 {
		t.Errorf("Expected 0 requests initially, got %d", count)
	}

	// Make some requests
	rl.Allow()
	rl.Allow()

	count = rl.GetRequestCount()
	if count != 2 {
		t.Errorf("Expected 2 requests after allowing 2, got %d", count)
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	rl := NewRateLimiter(1, 200*time.Millisecond)

	// Use the limit
	if !rl.Allow() {
		t.Error("First request should be allowed")
	}

	// Wait should block until a slot is available
	start := time.Now()
	rl.Wait()
	elapsed := time.Since(start)

	// Should have waited at least 180ms (allowing some tolerance)
	if elapsed < 180*time.Millisecond {
		t.Errorf("Wait should block for ~200ms, blocked for %v", elapsed)
	}
	if elapsed > 250*time.Millisecond {
		t.Errorf("Wait blocked too long: %v", elapsed)
	}
}

func TestRateLimiter_RequestsExpire(t *testing.T) {
	rl := NewRateLimiter(3, 100*time.Millisecond)

	// Fill up the limit
	for i := 0; i < 3; i++ {
		rl.Allow()
	}

	// Wait for requests to expire
	time.Sleep(150 * time.Millisecond)

	// All 3 slots should be available again
	allowed := 0
	for i := 0; i < 3; i++ {
		if rl.Allow() {
			allowed++
		}
	}

	if allowed != 3 {
		t.Errorf("Expected 3 requests to be allowed after expiry, got %d", allowed)
	}
}
