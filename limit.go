package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter to authorize the amount of requests per given time.
//
// To limit 10 requests per minute:
//
// rateLimiter := NewRateLimiter(10, 1*time.Minute)
// isAllowed := rateLimiter.Allow()
type RateLimiter struct {
	maxRequests   int
	leaseCycle    time.Duration
	requestsCount int
	// lock for lastReset.
	lock      sync.Mutex
	lastReset time.Time
}

// NewRateLimiter to create a new rate limiter with maxRequests and leaseCycle.
func NewRateLimiter(maxRequests int, leaseCycle time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		lastReset:   time.Now(),
		leaseCycle:  leaseCycle,
	}
}

// Allow checks if a request is allowed based on the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	// Check if we need to reset the counter.
	if time.Since(rl.lastReset) >= time.Minute {
		rl.requestsCount = 0
		rl.lastReset = time.Now()
	}

	// Check if we've reached the limit.
	if rl.requestsCount >= rl.maxRequests {
		return false
	}

	rl.requestsCount++
	return true
}

func main() {
	// Create a rate limiter that allows a maximum of 10 requests per minute
	limiter := NewRateLimiter(10)

	// Simulate requests
	for i := 0; i < 20; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d allowed\n", i)
		} else {
			fmt.Printf("Request %d denied\n", i)
		}
		time.Sleep(5 * time.Second) // Simulate request processing time
	}
}
