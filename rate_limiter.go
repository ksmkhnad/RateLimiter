package RateLimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	limits      map[string]limit
	timeWindows map[string][]time.Time
}

type limit struct {
	count    int
	duration time.Duration
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limits:      make(map[string]limit),
		timeWindows: make(map[string][]time.Time),
	}
}

func (r *RateLimiter) AddLimit(key string, count int, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limits[key] = limit{count, duration}
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	l, exists := r.limits[key]
	if !exists {
		return true
	}

	now := time.Now()
	window := r.timeWindows[key]

	newWindow := []time.Time{}
	for _, t := range window {
		if now.Sub(t) < l.duration {
			newWindow = append(newWindow, t)
		}
	}

	if len(newWindow) < l.count {
		newWindow = append(newWindow, now)
		r.timeWindows[key] = newWindow
		return true
	}

	r.timeWindows[key] = newWindow
	return false
}