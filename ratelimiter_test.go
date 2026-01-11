package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, 5)
	if rl.rate != 10 {
		t.Errorf("Expected rate 10, got %f", rl.rate)
	}
	if rl.burst != 5 {
		t.Errorf("Expected burst 5, got %d", rl.burst)
	}
	if rl.tokens != 5 {
		t.Errorf("Expected initial tokens 5, got %f", rl.tokens)
	}
}

func TestRateLimiter_Allow_BurstCapacity(t *testing.T) {
	rl := NewRateLimiter(1, 3) // 1 req/sec, burst of 3

	// Should allow burst requests immediately
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Errorf("Request %d should be allowed (burst)", i+1)
		}
	}

	// Next request should be denied (burst exhausted)
	if rl.Allow() {
		t.Error("Request should be denied after burst exhausted")
	}
}

func TestRateLimiter_Allow_TokenRefill(t *testing.T) {
	rl := NewRateLimiter(10, 1) // 10 req/sec, burst of 1

	// Use the initial token
	if !rl.Allow() {
		t.Error("First request should be allowed")
	}

	// Should be denied immediately
	if rl.Allow() {
		t.Error("Second request should be denied")
	}

	// Wait for token refill (100ms should give us 1 token at 10 req/sec)
	time.Sleep(150 * time.Millisecond)

	// Should be allowed now
	if !rl.Allow() {
		t.Error("Request should be allowed after token refill")
	}
}

func TestRateLimiter_Allow_RateLimit(t *testing.T) {
	rate := 5.0 // 5 requests per second
	rl := NewRateLimiter(rate, 10)

	// Exhaust initial burst
	for i := 0; i < 10; i++ {
		rl.Allow()
	}

	// Count allowed requests over 1 second
	start := time.Now()
	allowed := 0
	for time.Since(start) < 1*time.Second {
		if rl.Allow() {
			allowed++
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Should be approximately 5 requests allowed (with some tolerance)
	if allowed < 4 || allowed > 6 {
		t.Errorf("Expected ~5 requests allowed in 1 second, got %d", allowed)
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(100, 50) // 100 req/sec, burst of 50
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

	// Should allow at most burst capacity immediately
	if allowed > 50 {
		t.Errorf("Should not allow more than burst capacity, got %d", allowed)
	}
	if allowed < 40 {
		t.Errorf("Should allow close to burst capacity, got %d", allowed)
	}
}

func TestRateLimiter_GetTokens(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	// Initial tokens should be burst capacity
	tokens := rl.GetTokens()
	if tokens < 4.9 || tokens > 5.1 {
		t.Errorf("Expected ~5 tokens initially, got %f", tokens)
	}

	// Use some tokens
	rl.Allow()
	rl.Allow()

	tokens = rl.GetTokens()
	if tokens > 3.1 {
		t.Errorf("Expected ~3 tokens after using 2, got %f", tokens)
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	rl := NewRateLimiter(10, 1) // 10 req/sec, burst of 1

	// Use initial token
	if !rl.Allow() {
		t.Error("First request should be allowed")
	}

	// Wait should block until token is available
	start := time.Now()
	rl.Wait()
	elapsed := time.Since(start)

	// Should have waited at least 90ms (allowing some tolerance)
	if elapsed < 80*time.Millisecond {
		t.Errorf("Wait should block for ~100ms, blocked for %v", elapsed)
	}
}

func TestRateLimiter_MaxTokensCapped(t *testing.T) {
	rl := NewRateLimiter(10, 3) // 10 req/sec, burst of 3

	// Exhaust tokens
	for i := 0; i < 3; i++ {
		rl.Allow()
	}

	// Wait long enough to refill more than burst capacity
	time.Sleep(1 * time.Second)

	// Should only allow burst capacity, not more
	allowed := 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			allowed++
		}
	}

	if allowed != 3 {
		t.Errorf("Expected burst capacity (3) requests, got %d", allowed)
	}
}
