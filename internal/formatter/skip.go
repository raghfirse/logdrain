package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Skipper drops the first N log lines (or first N per key) before passing
// subsequent lines through. It is the complement of Head.
type Skipper struct {
	limit  int
	keyField string
	counts map[string]int
}

// NewSkip creates a Skipper that skips the first n lines.
// If keyField is non-empty, the skip count is tracked independently per
// distinct value of that field in the JSON object.
func NewSkip(n int, keyField string) (*Skipper, error) {
	if n < 0 {
		return nil, fmt.Errorf("skip: limit must be >= 0, got %d", n)
	}
	return &Skipper{
		limit:    n,
		keyField: keyField,
		counts:   make(map[string]int),
	}, nil
}

// Apply returns (line, true) once the skip threshold has been reached,
// and ("", false) for lines that are still being skipped.
func (s *Skipper) Apply(line string) (string, bool) {
	if s.limit == 0 {
		return line, true
	}

	key := ""
	if s.keyField != "" {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &obj); err == nil {
			if raw, ok := findRawValueCaseInsensitive(obj, s.keyField); ok {
				key = string(raw)
			}
		}
	}

	s.counts[key]++
	if s.counts[key] <= s.limit {
		return "", false
	}
	return line, true
}

// Reset clears all counters.
func (s *Skipper) Reset() {
	s.counts = make(map[string]int)
}

// ParseSkipFlag parses a skip flag value of the form "N" or "N:field".
func ParseSkipFlag(val string) (*Skipper, error) {
	if val == "" {
		return NewSkip(0, "")
	}
	for i, c := range val {
		if c == ':' {
			n, err := strconv.Atoi(val[:i])
			if err != nil {
				return nil, fmt.Errorf("skip: invalid number %q", val[:i])
			}
			return NewSkip(n, val[i+1:])
		}
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("skip: invalid value %q, want N or N:field", val)
	}
	return NewSkip(n, "")
}
