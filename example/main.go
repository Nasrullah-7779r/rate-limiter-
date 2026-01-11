package main

import (
	"fmt"
	"time"

	ratelimiter "github.com/Nasrullah-7779r/rate-limiter-"
)

func main() {
	// Example 1: Basic usage with sliding window
	fmt.Println("Example 1: Basic Sliding Window Log Rate Limiting")
	fmt.Println("Creating rate limiter: 5 requests per second")
	rl := ratelimiter.NewRateLimiter(5, time.Second)

	fmt.Println("\nSending 10 requests immediately:")
	for i := 1; i <= 10; i++ {
		if rl.Allow() {
			fmt.Printf("Request %d: ✓ Allowed\n", i)
		} else {
			fmt.Printf("Request %d: ✗ Denied\n", i)
		}
	}

	// Example 2: Rate limiting over time with window sliding
	fmt.Println("\n\nExample 2: Sliding Window Over Time")
	fmt.Println("Creating rate limiter: 3 requests per 500ms")
	rl2 := ratelimiter.NewRateLimiter(3, 500*time.Millisecond)

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
	fmt.Printf("\nTotal requests allowed: %d\n", count)

	// Example 3: Using Wait method
	fmt.Println("\n\nExample 3: Using Wait Method")
	fmt.Println("Creating rate limiter: 2 requests per 500ms")
	rl3 := ratelimiter.NewRateLimiter(2, 500*time.Millisecond)

	fmt.Println("\nSending 5 requests with Wait (blocks until allowed):")
	for i := 1; i <= 5; i++ {
		start := time.Now()
		rl3.Wait()
		elapsed := time.Since(start)
		fmt.Printf("Request %d: processed after waiting %v\n", i, elapsed.Round(time.Millisecond))
	}

	// Example 4: Monitoring request count
	fmt.Println("\n\nExample 4: Monitoring Request Count")
	rl4 := ratelimiter.NewRateLimiter(5, time.Second)

	fmt.Printf("Initial request count: %d\n", rl4.GetRequestCount())

	for i := 1; i <= 3; i++ {
		rl4.Allow()
		fmt.Printf("After request %d, count in window: %d\n", i, rl4.GetRequestCount())
	}

	time.Sleep(1100 * time.Millisecond)
	fmt.Printf("After 1.1s wait (window expired), count: %d\n", rl4.GetRequestCount())
}
