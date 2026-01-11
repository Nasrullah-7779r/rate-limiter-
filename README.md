# Rate Limiter

A simple and efficient rate limiter implementation in Go using the token bucket algorithm.

## Features

- **Token Bucket Algorithm**: Efficient and widely-used rate limiting algorithm
- **Thread-Safe**: Safe for concurrent use across multiple goroutines
- **Configurable**: Set custom rate limits and burst capacity
- **Flexible**: Support for both immediate checks and blocking waits
- **Zero Dependencies**: Built using only Go standard library

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
    ratelimiter "github.com/Nasrullah-7779r/rate-limiter-"
)

func main() {
    // Create a rate limiter: 10 requests per second, burst of 5
    rl := ratelimiter.NewRateLimiter(10, 5)
    
    // Check if a request is allowed
    if rl.Allow() {
        fmt.Println("Request allowed")
    } else {
        fmt.Println("Request denied - rate limit exceeded")
    }
}
```

### Wait for Available Token

```go
// Block until a request is allowed
rl.Wait()
fmt.Println("Request processed")
```

### Monitor Available Tokens

```go
tokens := rl.GetTokens()
fmt.Printf("Available tokens: %.2f\n", tokens)
```

## API Reference

### `NewRateLimiter(rate float64, burst int) *RateLimiter`

Creates a new rate limiter with the specified rate and burst capacity.

- `rate`: Number of requests allowed per second
- `burst`: Maximum number of requests that can be made immediately (bucket capacity)

### `Allow() bool`

Checks if a request is allowed based on the current rate limit. Returns `true` if allowed, `false` otherwise.

### `Wait()`

Blocks until a request is allowed. Use this when you want to ensure the request is processed rather than being denied.

### `GetTokens() float64`

Returns the current number of available tokens. Useful for monitoring and debugging.

## How It Works

The rate limiter uses the **token bucket algorithm**:

1. Tokens are added to the bucket at a constant rate (the configured rate)
2. The bucket has a maximum capacity (the burst capacity)
3. Each request consumes one token
4. If tokens are available, the request is allowed
5. If no tokens are available, the request is denied (or waits)

This allows for:
- Smooth rate limiting over time
- Burst handling for temporary spikes
- Efficient token refill without timers

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