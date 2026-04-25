package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rollup groups log lines by a key field and emits a summary line
// containing the group key and count of lines seen in that group.
type Rollup struct {
	field  string
	counts map[string]int
}

// NewRollup creates a Rollup that groups by the given field name.
// Returns an error if field is empty.
func NewRollup(field string) (*Rollup, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("rollup: field name must not be empty")
	}
	return &Rollup{
		field:  field,
		counts: make(map[string]int),
	}, nil
}

// Apply accumulates the line into the rollup bucket for its field value.
// Returns an empty string (suppresses the line) until Flush is called.
func (r *Rollup) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	key := findRawValueCaseInsensitive(obj, r.field)
	if key == "" {
		key = "(missing)"
	}
	r.counts[key]++
	return ""
}

// Flush emits one summary JSON line per unique key value, then resets state.
func (r *Rollup) Flush() []string {
	var lines []string
	for k, n := range r.counts {
		b, _ := json.Marshal(map[string]interface{}{
			r.field: k,
			"count": n,
		})
		lines = append(lines, string(b))
	}
	r.counts = make(map[string]int)
	return lines
}

// ParseRollupFlag parses a rollup field name from a flag string.
func ParseRollupFlag(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("rollup: field name must not be empty")
	}
	return s, nil
}
