package main

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	t.Run("Allow within limit", func(t *testing.T) {
		rateLimiter := NewRateLimiter(10, 1*time.Second)
		for i := 0; i < 5; i++ {
			if !rateLimiter.Allow() {
				t.Errorf("Request %d should be allowed", i)
			}
		}
	})

	t.Run("Deny exceeding limit", func(t *testing.T) {
		rateLimiter := NewRateLimiter(5, 1*time.Second)
		for i := 0; i < 10; i++ {
			if i < 5 {
				if !rateLimiter.Allow() {
					t.Errorf("Request %d should be allowed", i)
				}
			} else {
				if rateLimiter.Allow() {
					t.Errorf("Request %d should be denied", i)
				}
			}
		}
	})

	t.Run("Reset counter after lease cycle", func(t *testing.T) {
		rateLimiter := NewRateLimiter(5, 1*time.Second)
		for i := 0; i < 10; i++ {
			if i < 5 {
				if !rateLimiter.Allow() {
					t.Errorf("Request %d should be allowed", i)
				}
			} else {
				if rateLimiter.Allow() {
					t.Errorf("Request %d should be denied", i)
				}
			}
		}
		time.Sleep(1 * time.Second)
		for i := 0; i < 5; i++ {
			if !rateLimiter.Allow() {
				t.Errorf("Request %d should be allowed after reset", i)
			}
		}
	})
}
