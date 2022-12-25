package ratelimiter

import (
	"sync"
	"time"
)

type bucket struct {
	mu       sync.Mutex
	doneChan chan struct{}
	tokens   int
	rate     time.Duration
}

func NewBucket(rate time.Duration) *bucket {
	return &bucket{
		tokens:   1,
		doneChan: make(chan struct{}),
		rate:     rate,
	}
}

func (b *bucket) Start() {
	ticker := time.NewTicker(b.rate)

	for {
		select {
		case <-b.doneChan:
			return
		case <-ticker.C:
			b.mu.Lock()
			if b.tokens < 1 {
				b.tokens++
			}
			b.mu.Unlock()
		}
	}
}

func (b *bucket) RequestToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}
