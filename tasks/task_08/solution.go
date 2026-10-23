package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	var t time.Time
	if clock != nil {
		t = clock.Now()
	}

	return &Limiter {
		clock: clock,
		rate: ratePerSec,
		burst: burst,
		tokens: float64(burst),
		last: t,
	}
}
func (l *Limiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if n <= 0 || n > l.burst{
		return false 
	}

	if l.burst <=0 || l.clock == nil {
		return false
	}

	now := l.clock.Now()
	if now.Before(l.last) {
		now = l.last
	}

	if l.rate > 0 {
		sec := now.Sub(l.last).Seconds()
		l.tokens += sec * l.rate
		if l.tokens > float64(l.burst) {
			l.tokens = float64(l.burst)
		}
		l.last = now
	}

	if l.tokens >= float64(n) {
		l.tokens -= float64(n)
		return true
	}

	return false
}