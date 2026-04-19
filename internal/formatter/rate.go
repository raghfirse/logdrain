package formatter

import (
	"sync"
	"time"
)

// RateLimiter suppresses log lines that exceed a threshold within a window.
type RateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	max      int
	counts   map[string][]time.Time
	nowFn    func() time.Time
}

// NewRateLimiter creates a RateLimiter allowing at most max lines per key per window.
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		window: window,
		max:    max,
		counts: make(map[string][]time.Time),
		nowFn:  time.Now,
	}
}

// Allow returns true if the line identified by key is within the rate limit.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.nowFn()
	cutoff := now.Add(-r.window)

	times := r.counts[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= r.max {
		r.counts[key] = filtered
		return false
	}

	r.counts[key] = append(filtered, now)
	return true
}

// Reset clears all rate limit state.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counts = make(map[string][]time.Time)
}
