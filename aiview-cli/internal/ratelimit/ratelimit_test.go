package ratelimit

import (
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	l := New(10, 5) // 10 tokens/sec, burst 5

	// Should allow 5 requests immediately
	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Errorf("expected Allow() to return true for request %d", i)
		}
	}

	// 6th should fail (burst exhausted)
	if l.Allow() {
		t.Error("expected Allow() to return false after burst exhausted")
	}
}

func TestLimiter_Wait(t *testing.T) {
	l := New(100, 1) // 100 tokens/sec, burst 1

	start := time.Now()
	for i := 0; i < 3; i++ {
		l.Wait()
	}
	elapsed := time.Since(start)

	// Should take at least 10ms for 2 waits at 100/sec
	if elapsed < 10*time.Millisecond {
		t.Errorf("expected at least 10ms, got %v", elapsed)
	}
}
