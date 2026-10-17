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

	l := Limiter{
		clock:  clock,
		burst:  burst,
		tokens: float64(burst),
		rate:   ratePerSec,
	}

	if clock != nil {
		l.last = l.clock.Now()
	}

	return &l
}

func (l *Limiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.burst <= 0 || l.clock == nil {
		return false
	}

	if n > l.burst || n <= 0 {
		return false
	}

	//время может измениться, пока выполняются проверки
	nowClock := l.clock.Now()
	if nowClock.After(l.last) && l.rate > 0 {
		r := nowClock.Sub(l.last).Seconds() * l.rate
		l.tokens = min(l.tokens+r, float64(l.burst))
		l.last = nowClock
	}

	if l.tokens >= float64(n) {
		l.tokens -= float64(n)
		return true
	}
	return false
}
