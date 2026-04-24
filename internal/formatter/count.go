package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Counter tracks the number of occurrences of each distinct value for a given
// JSON field and emits a summary line when Flush is called.
type Counter struct {
	field  string
	counts map[string]int64
}

// NewCounter creates a Counter that groups log lines by the value of field.
// Returns an error if field is empty.
func NewCounter(field string) (*Counter, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("count: field name must not be empty")
	}
	return &Counter{
		field:  field,
		counts: make(map[string]int64),
	}, nil
}

// Apply records the value of the configured field from line.
// It always returns an empty string so that individual lines are suppressed;
// call Flush to obtain the accumulated summary.
func (c *Counter) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return ""
	}
	raw, ok := findRawValueCaseInsensitive(obj, c.field)
	if !ok {
		return ""
	}
	var key string
	if err := json.Unmarshal(raw, &key); err != nil {
		key = string(raw)
	}
	c.counts[key]++
	return ""
}

// Flush returns a JSON line for each distinct value observed, sorted by count
// descending, and resets internal state. Returns nil if no data was recorded.
func (c *Counter) Flush() []string {
	if len(c.counts) == 0 {
		return nil
	}
	results := make([]string, 0, len(c.counts))
	for val, n := range c.counts {
		b, _ := json.Marshal(map[string]interface{}{
			c.field: val,
			"count": n,
		})
		results = append(results, string(b))
	}
	c.counts = make(map[string]int64)
	return results
}

// ParseCountFlag parses a bare field name from the --count flag value.
func ParseCountFlag(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("count: expected a field name, got empty string")
	}
	return s, nil
}
