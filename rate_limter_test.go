package RateLimiter

import (
	"testing"
	"time"
)

func TestTokenBucket_Allow_AllowedWhenUnderLimit(t *testing.T) {
	tb := NewTokenBucket("test", 5, time.Second)
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("Expected request %d to be allowed, but it was not", i+1)
		}
	}
}

func TestTokenBucket_Allow_NotAllowedWhenOverLimit(t *testing.T) {
	tb := NewTokenBucket("test", 5, time.Second)
	for i := 0; i < 5; i++ {
		tb.Allow()
	}
	if tb.Allow() {
		t.Error("Expected request to be not allowed, but it was allowed")
	}
}

func TestTokenBucket_Allow_AllowedAfterRefill(t *testing.T) {
	tb := NewTokenBucket("test", 1, time.Second)
	tb.Allow()
	time.Sleep(2 * time.Second)
	if !tb.Allow() {
		t.Error("Expected request to be allowed after refill, but it was not")
	}
}
