package formatter

import (
	"sync"
	"time"
)

// DedupeFilter suppresses repeated identical log lines within a time window.
type DedupeFilter struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     func() time.Time
}

// NewDedupeFilter creates a DedupeFilter with the given suppression window.
func NewDedupeFilter(window time.Duration) *DedupeFilter {
	return &DedupeFilter{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// IsDuplicate returns true if the line was seen within the window.
func (d *DedupeFilter) IsDuplicate(line string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	d.evict(now)

	if _, ok := d.seen[line]; ok {
		return true
	}
	d.seen[line] = now
	return false
}

// evict removes entries older than the window. Must be called with lock held.
func (d *DedupeFilter) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) > d.window {
			delete(d.seen, k)
		}
	}
}

// Reset clears all seen entries.
func (d *DedupeFilter) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}
