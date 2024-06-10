package RateLimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	limits  map[string]limit
	windows map[string]*ringBuffer
}

type limit struct {
	count    int
	duration time.Duration
}

type ringBuffer struct {
	buffer []time.Time
	size   int
	start  int
	end    int
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		buffer: make([]time.Time, size),
		size:   size,
	}
}

func (rb *ringBuffer) add(t time.Time) {
	rb.buffer[rb.end] = t
	rb.end = (rb.end + 1) % rb.size
	if rb.end == rb.start {
		rb.start = (rb.start + 1) % rb.size
	}
}

func (rb *ringBuffer) countWithin(duration time.Duration) int {
	count := 0
	now := time.Now()
	for i := 0; i < rb.size; i++ {
		index := (rb.start + i) % rb.size
		if now.Sub(rb.buffer[index]) < duration {
			count++
		}
	}
	return count
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limits:  make(map[string]limit),
		windows: make(map[string]*ringBuffer),
	}
}

func (r *RateLimiter) AddLimit(key string, count int, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limits[key] = limit{count, duration}
	r.windows[key] = newRingBuffer(count)
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	l, exists := r.limits[key]
	if !exists {
		return true
	}

	window, exists := r.windows[key]
	if !exists {
		return true
	}

	now := time.Now()
	elapsed := now.Sub(window.buffer[window.start])

	if elapsed >= l.duration {
		for i := range window.buffer {
			window.buffer[i] = now
		}
		window.start = 0
		window.end = 0
		return true
	}

	if window.countWithin(l.duration) < l.count {
		window.add(now)
		return true
	}

	return false
}

