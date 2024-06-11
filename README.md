# Rate Limiter

This Go package provides a simple rate limiting mechanism to control the frequency of specific tasks. It allows you to define rate limits and check if a task is allowed based on those limits.

## Installation

```sh
go get github.com/ksmkhnad/RateLimiter

## Usage example
## main.go
package main

import (
	"time"

	ratelimiter "github.com/ksmkhnad/RateLimiter"
)

func main() {
	// Пользователь может отправлять не более 5 сообщений в секунду
	userLimiter := ratelimiter.NewTokenBucket("user", 5, time.Second)

	// Один IP-адрес может отправлять не более 10000 запросов в минуту
	ipLimiter := ratelimiter.NewTokenBucket("ip", 10000, time.Minute)

	// Юзер может иметь не более 3-х неудачных транзакций по карте в день
	cardLimiter := ratelimiter.NewTokenBucket("card", 3, 24*time.Hour)

	for i := 0; i < 10; i++ {
		if userLimiter.Allow() {
			println("User message allowed")
		} else {
			println("User message not allowed")
		}
		time.Sleep(200 * time.Millisecond)
	}

	for i := 0; i < 20000; i++ {
		if ipLimiter.Allow() {
			println("IP request allowed")
		} else {
			println("IP request not allowed")
		}
		time.Sleep(3 * time.Millisecond)
	}
	
	for i := 0; i < 5; i++ {
		if cardLimiter.Allow() {
			println("Card transaction allowed")
		} else {
			println("Card transaction not allowed")
		}
		time.Sleep(10 * time.Second)
	}
}
