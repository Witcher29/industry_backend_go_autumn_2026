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
	limiterPtr := &Limiter{
		clock: clock,
		rate:  ratePerSec,
		burst: burst,
	}
	if clock == nil || burst <= 0 {
		return limiterPtr
	}
	limiterPtr.tokens = float64(burst)
	limiterPtr.last = clock.Now()
	return limiterPtr
}

func (limiterPtr *Limiter) AllowN(n int) bool {
	limiterPtr.mu.Lock()
	defer limiterPtr.mu.Unlock()

	if limiterPtr.clock == nil || limiterPtr.burst <= 0 || n <= 0 || n > limiterPtr.burst {
		return false
	}

	now := limiterPtr.clock.Now()

	if now.After(limiterPtr.last) {
		if limiterPtr.rate > 0 {
			elapsedSeconds := now.Sub(limiterPtr.last).Seconds()
			limiterPtr.tokens += elapsedSeconds * limiterPtr.rate
			if limiterPtr.tokens > float64(limiterPtr.burst) {
				limiterPtr.tokens = float64(limiterPtr.burst)
			}
		}
		limiterPtr.last = now
	}

	if limiterPtr.tokens < float64(n) {
		return false
	}
	limiterPtr.tokens -= float64(n)
	return true
}
