package main

import (
	"fmt"
	"time"

	ratelimiter "github.com/Nasrullah-7779r/rate-limiter-"
)

func main() {
	// Example 1: Basic usage with burst capacity
	fmt.Println("Example 1: Basic Rate Limiting")
	fmt.Println("Creating rate limiter: 2 requests/second, burst of 5")
	rl := ratelimiter.NewRateLimiter(2, 5)

	fmt.Println("\nSending 10 requests immediately:")
	for i := 1; i <= 10; i++ {
		if rl.Allow() {
			fmt.Printf("Request %d: ✓ Allowed\n", i)
		} else {
			fmt.Printf("Request %d: ✗ Denied\n", i)
		}
	}

	// Example 2: Rate limiting over time
	fmt.Println("\n\nExample 2: Rate Limiting Over Time")
	fmt.Println("Creating rate limiter: 5 requests/second, burst of 2")
	rl2 := ratelimiter.NewRateLimiter(5, 2)

	fmt.Println("\nSending requests over 2 seconds:")
	start := time.Now()
	count := 0
	for time.Since(start) < 2*time.Second {
		if rl2.Allow() {
			count++
			fmt.Printf("Request %d allowed at %v\n", count, time.Since(start).Round(time.Millisecond))
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\nTotal requests allowed: %d (expected ~10-12)\n", count)

	// Example 3: Using Wait method
	fmt.Println("\n\nExample 3: Using Wait Method")
	fmt.Println("Creating rate limiter: 3 requests/second, burst of 1")
	rl3 := ratelimiter.NewRateLimiter(3, 1)

	fmt.Println("\nSending 5 requests with Wait (blocks until allowed):")
	for i := 1; i <= 5; i++ {
		start := time.Now()
		rl3.Wait()
		elapsed := time.Since(start)
		fmt.Printf("Request %d: processed after waiting %v\n", i, elapsed.Round(time.Millisecond))
	}

	// Example 4: Monitoring available tokens
	fmt.Println("\n\nExample 4: Monitoring Available Tokens")
	rl4 := ratelimiter.NewRateLimiter(10, 5)

	fmt.Printf("Initial tokens: %.2f\n", rl4.GetTokens())

	for i := 1; i <= 3; i++ {
		rl4.Allow()
		fmt.Printf("After request %d, tokens: %.2f\n", i, rl4.GetTokens())
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("After 200ms wait, tokens: %.2f\n", rl4.GetTokens())
}
