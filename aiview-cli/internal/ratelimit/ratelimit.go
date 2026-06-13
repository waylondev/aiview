// Package ratelimit provides a token bucket rate limiter.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter implements a token bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	rate     float64   // tokens per second
	burst    int       // max tokens
	tokens   float64
	lastTime time.Time
}

// New creates a new Limiter with the given rate (tokens/sec) and burst size.
func New(rate float64, burst int) *Limiter {
	return &Limiter{
		rate:     rate,
		burst:    burst,
		tokens:   float64(burst),
		lastTime: time.Now(),
	}
}

// Wait blocks until a token is available.
func (l *Limiter) Wait() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTime).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastTime = now

	if l.tokens < 1 {
		waitDuration := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))
		l.mu.Unlock()
		time.Sleep(waitDuration)
		l.mu.Lock()
		l.tokens = 0
	} else {
		l.tokens--
	}
}

// Allow returns true if a token is available without waiting.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTime).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastTime = now

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}
