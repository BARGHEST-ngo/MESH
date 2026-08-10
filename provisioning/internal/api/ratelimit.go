package api

import (
	"sync"

	"golang.org/x/time/rate"
)

type rateLimiter struct {
	mu         sync.Mutex
	refillRate rate.Limit
	burstSize  int
	limits     map[string]*rate.Limiter
}

func NewLimiter(r rate.Limit, burst int) *rateLimiter {
	return &rateLimiter{
		refillRate: r,
		burstSize:  burst,
		limits:     make(map[string]*rate.Limiter),
	}
}

func (l *rateLimiter) allow(keyID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	keyLimiter, ok := l.limits[keyID]
	if !ok {
		l.limits[keyID] = rate.NewLimiter(l.refillRate, l.burstSize)
		keyLimiter = l.limits[keyID]
	}
	return keyLimiter.Allow()
}
