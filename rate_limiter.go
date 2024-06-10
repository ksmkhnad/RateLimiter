package RateLimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	limits  map[string]limit
	windows map[string]*timeWindow
}

type limit struct {
	count    int
	duration time.Duration
}

type timeWindow struct {
	times []time.Time
	start int
	end   int
	size  int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limits:  make(map[string]limit),
		windows: make(map[string]*timeWindow),
	}
}

func (r *RateLimiter) AddLimit(key string, count int, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limits[key] = limit{count, duration}
	r.windows[key] = &timeWindow{
		times: make([]time.Time, count),
		size:  count,
	}
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	lim, exists := r.limits[key]
	window, winExists := r.windows[key]
	r.mu.Unlock()

	if !exists || !winExists {
		return true
	}

	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	// Remove expired timestamps
	for i := 0; i < window.size; i++ {
		if now.Sub(window.times[window.start]) < lim.duration {
			break
		}
		window.start = (window.start + 1) % window.size
		window.end = (window.end + 1) % window.size
	}

	// Check if we can allow a new request
	if (window.end+1)%window.size != window.start {
		window.times[window.end] = now
		window.end = (window.end + 1) % window.size
		return true
	}

	return false
}
