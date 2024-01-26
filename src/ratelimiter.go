package main

import (
	"time"
)

// Sliding window rate limiter allowing a certain number of events to occur in a specified time window

type RateLimiter struct {
	Events       []time.Time
	TimeInterval time.Duration
	MaxEvents    int
}

func NewRateLimiter(interval time.Duration, max int) *RateLimiter {
	return &RateLimiter{
		TimeInterval: interval,
		MaxEvents:    max,
	}
}

func (r *RateLimiter) IsAllowed() bool {
	now := time.Now()

	// Remove timestamps outside the current window
	for len(r.Events) > 0 && now.Sub(r.Events[0]) > r.TimeInterval {
		r.Events = r.Events[1:]
	}

	if len(r.Events) < r.MaxEvents {
		r.Events = append(r.Events, now)
		return true
	}

	return false
}
