package RateLimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	key      string
	tokens   int
	interval time.Duration
	bucket   chan struct{}
	mu       sync.Mutex
}

func NewTokenBucket(key string, tokens int, interval time.Duration) *TokenBucket {
	tb := &TokenBucket{
		key:      key,
		tokens:   tokens,
		interval: interval,
		bucket:   make(chan struct{}, tokens),
	}

	for i := 0; i < tokens; i++ {
		tb.bucket <- struct{}{}
	}

	go tb.refill()

	return tb
}

func (tb *TokenBucket) refill() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		tb.mu.Lock()
		if len(tb.bucket) < tb.tokens {
			tb.bucket <- struct{}{}
		}
		tb.mu.Unlock()
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	select {
	case <-tb.bucket:
		return true
	default:
		return false
	}
}
