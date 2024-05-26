package main

import (
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
	if time.Since(rl.lastReset) >= rl.leaseCycle {
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

// ConsistData to save the rate limit data on distributed system.
func (rl *RateLimiter) ConsistData() {}
