package formatter

import (
	"fmt"
	"strconv"
	"strings"
)

// Head limits output to the first N lines per unique key (or globally if key is empty).
type Head struct {
	key   string
	limit int
	counts map[string]int
}

// NewHead creates a Head limiter. limit is the max number of lines to pass through.
// key is an optional JSON field name; if non-empty, the limit is applied per distinct value of that field.
func NewHead(limit int, key string) (*Head, error) {
	if limit < 0 {
		return nil, fmt.Errorf("head: limit must be non-negative, got %d", limit)
	}
	return &Head{
		key:    strings.ToLower(key),
		limit:  limit,
		counts: make(map[string]int),
	}, nil
}

// Allow returns true if the line should be passed through.
func (h *Head) Allow(line string) bool {
	if h.limit == 0 {
		return true
	}
	bucket := "_global_"
	if h.key != "" {
		bucket = findKeyCaseInsensitive(line, h.key)
	}
	h.counts[bucket]++
	return h.counts[bucket] <= h.limit
}

// Reset clears all counters.
func (h *Head) Reset() {
	h.counts = make(map[string]int)
}

// ParseHeadFlag parses a head flag value of the form "N" or "N:field".
func ParseHeadFlag(s string) (*Head, error) {
	if s == "" {
		return nil, fmt.Errorf("head: flag value must not be empty")
	}
	parts := strings.SplitN(s, ":", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("head: invalid limit %q: %w", parts[0], err)
	}
	key := ""
	if len(parts) == 2 {
		key = parts[1]
	}
	return NewHead(n, key)
}
