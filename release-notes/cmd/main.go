package main

import (
	"time"
)

type (
	RateLimiter struct {
		tokens          chan Token
		refreshInterval time.Duration
	}

	Token struct{}
)

func (r *RateLimiter) Take() {
	<-r.tokens

	// LOOK HERE
	go func() {
		time.Sleep(r.refreshInterval)
		r.tokens <- Token{}
	}()
}

func New(size int, refreshInterval time.Duration) *RateLimiter {
	tokens := make(chan Token, size)

	for range size {
		tokens <- Token{}
	}

	return &RateLimiter{
		tokens:          tokens,
		refreshInterval: refreshInterval,
	}
}
