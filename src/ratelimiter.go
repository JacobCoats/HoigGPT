package main

import (
	"time"
)

type RateLimiter struct {
	LastRecordedTime time.Time
	TimeInterval     time.Duration
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		TimeInterval: interval,
	}
}

func (r *RateLimiter) IsAllowed() bool {
	now := time.Now()
	if now.Sub(r.LastRecordedTime) >= r.TimeInterval {
		r.LastRecordedTime = now
		return true
	}
	return false
}
