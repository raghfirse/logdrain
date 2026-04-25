package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// TopEntry holds a key value and its count.
type TopEntry struct {
	Value string
	Count int
}

// Top tracks the most frequent values for a given JSON field.
type Top struct {
	field  string
	n      int
	counts map[string]int
}

// NewTop creates a Top tracker for the given field, returning the top n values on Flush.
func NewTop(field string, n int) (*Top, error) {
	if strings.TrimSpace(field) == "" {
		return nil, fmt.Errorf("top: field must not be empty")
	}
	if n <= 0 {
		return nil, fmt.Errorf("top: n must be greater than zero, got %d", n)
	}
	return &Top{field: field, n: n, counts: make(map[string]int)}, nil
}

// Apply records the field value from the line and suppresses output (returns "", false).
func (t *Top) Apply(line string) (string, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true
	}
	raw := findRawValueCaseInsensitive(obj, t.field)
	if raw == nil {
		return "", false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		s = strings.Trim(string(raw), `"`)
	}
	t.counts[s]++
	return "", false
}

// Flush returns the top-n entries sorted by count descending as JSON lines.
func (t *Top) Flush() []string {
	entries := make([]TopEntry, 0, len(t.counts))
	for v, c := range t.counts {
		entries = append(entries, TopEntry{Value: v, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	if len(entries) > t.n {
		entries = entries[:t.n]
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		b, _ := json.Marshal(map[string]interface{}{t.field: e.Value, "count": e.Count})
		out = append(out, string(b))
	}
	return out
}

// Reset clears accumulated counts.
func (t *Top) Reset() {
	t.counts = make(map[string]int)
}

// ParseTopFlag parses "field:n" into (field, n, error).
func ParseTopFlag(s string) (string, int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("top: expected field:n, got %q", s)
	}
	var n int
	if _, err := fmt.Sscanf(parts[1], "%d", &n); err != nil {
		return "", 0, fmt.Errorf("top: invalid n %q: %w", parts[1], err)
	}
	return parts[0], n, nil
}
