# Rate Limiter

This Go package provides a simple rate limiting mechanism to control the frequency of specific tasks. It allows you to define rate limits and check if a task is allowed based on those limits.

## Installation

```sh
go get github.com/ksmkhnad/RateLimiter

## Usage example
## main.go
package main

import (
	"fmt"
	"time"

	ratelimiter "github.com/ksmkhnad/RateLimiter"
)

func main() {
	rl := ratelimiter.NewRateLimiter()
	rl.AddLimit("login:user:123", 5, time.Second)
	rl.AddLimit("login:ip:192.168.1.1", 10000, time.Minute)
	rl.AddLimit("txn:card:123", 3, 24*time.Hour)

	for i := 0; i < 10; i++ {
		if rl.Allow("login:user:123") {
			fmt.Println("User login allowed")
		} else {
			fmt.Println("User login rate limit exceeded")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
