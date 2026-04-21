package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Limiter stops emitting lines after a global count is reached.
// Unlike Head, Limiter operates purely on line count with no per-key grouping.
type Limiter struct {
	max   int
	count int
}

// NewLimiter creates a Limiter that allows at most max lines through.
// A max of 0 disables limiting.
func NewLimiter(max int) (*Limiter, error) {
	if max < 0 {
		return nil, fmt.Errorf("limit: max must be non-negative, got %d", max)
	}
	return &Limiter{max: max}, nil
}

// Apply returns the line unchanged if under the limit, or ("" , true) once
// the limit is reached. The boolean signals that the stream should stop.
func (l *Limiter) Apply(line string) (string, bool) {
	if l.max == 0 {
		return line, false
	}
	if l.count >= l.max {
		return "", true
	}
	l.count++
	return line, false
}

// Reset clears the internal counter.
func (l *Limiter) Reset() {
	l.count = 0
}

// ParseLimitFlag parses a string into a non-negative integer limit.
func ParseLimitFlag(s string) (int, error) {
	if s == "" || s == "0" {
		return 0, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("limit: invalid value %q: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("limit: value must be non-negative, got %d", n)
	}
	return n, nil
}

// limitInfo is used internally to pretty-print limit state.
type limitInfo struct {
	Max   int `json:"max"`
	Count int `json:"count"`
}

func (l *Limiter) String() string {
	b, _ := json.Marshal(limitInfo{Max: l.max, Count: l.count})
	return string(b)
}
