# Rate Limiter - Sliding Window Log

A simple and efficient rate limiter implementation in Go using the Sliding Window Log algorithm.

## Features

- **Sliding Window Log Algorithm**: Precise rate limiting with accurate request tracking
- **Thread-Safe**: Safe for concurrent use across multiple goroutines
- **Configurable**: Set custom request limits and time windows
- **Flexible**: Support for both immediate checks and blocking waits
- **Zero Dependencies**: Built using only Go standard library
- **Memory Efficient**: Automatically cleans up expired request logs

## Installation

```bash
go get github.com/Nasrullah-7779r/rate-limiter-
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    ratelimiter "github.com/Nasrullah-7779r/rate-limiter-"
)

func main() {
    // Create a rate limiter: 10 requests per second
    rl := ratelimiter.NewRateLimiter(10, time.Second)

    // Check if a request is allowed
    if rl.Allow() {
        fmt.Println("Request allowed")
    } else {
        fmt.Println("Request denied - rate limit exceeded")
    }
}
```

### Wait for Available Slot

```go
// Block until a request is allowed
rl.Wait()
fmt.Println("Request processed")
```

### Monitor Request Count

```go
count := rl.GetRequestCount()
fmt.Printf("Requests in current window: %d\n", count)
```

## API Reference

### `NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter`

Creates a new rate limiter with the specified maximum requests and time window.

- `maxRequests`: Maximum number of requests allowed within the time window
- `window`: Time window duration (e.g., `time.Second`, `time.Minute`)

### `Allow() bool`

Checks if a request is allowed based on the current rate limit. Returns `true` if allowed, `false` otherwise.

### `Wait()`

Blocks until a request is allowed. Use this when you want to ensure the request is processed rather than being denied.

### `GetRequestCount() int`

Returns the current number of requests in the sliding window. Useful for monitoring and debugging.

## How It Works

The rate limiter uses the **Sliding Window Log algorithm**:

1. Each request timestamp is logged in memory
2. When a new request arrives, old requests outside the time window are removed
3. If the number of requests within the window is less than the limit, the request is allowed
4. Otherwise, the request is denied (or waits until the oldest request expires)

This provides:

- **Precise rate limiting**: Exact tracking of requests within the time window
- **No burst issues**: Prevents request spikes at window boundaries
- **Fair distribution**: Requests are smoothly distributed over time
- **Predictable behavior**: Easy to understand and reason about

## Example Scenarios

### HTTP API Rate Limiting

```go
// Allow 100 requests per minute
rl := ratelimiter.NewRateLimiter(100, time.Minute)
```

### Database Query Throttling

```go
// Allow 10 queries per second
rl := ratelimiter.NewRateLimiter(10, time.Second)
```

### User Action Limiting

```go
// Allow 5 login attempts per 5 minutes
rl := ratelimiter.NewRateLimiter(5, 5*time.Minute)
```

## Testing

Run the tests:

```bash
go test -v
```

Run the example:

```bash
cd example
go run main.go
```

## Use Cases

- API rate limiting
- Request throttling
- Resource usage control
- DoS protection
- Fair resource allocation

## License

MIT License
