package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Sorter buffers log lines and emits them sorted by a numeric or string field.
type Sorter struct {
	field     string
	descending bool
	buffer    []string
}

// NewSorter creates a Sorter that orders buffered lines by field.
// direction must be "asc" or "desc".
func NewSorter(field, direction string) (*Sorter, error) {
	if field == "" {
		return nil, fmt.Errorf("sort: field must not be empty")
	}
	desc := false
	switch strings.ToLower(direction) {
	case "desc", "descending":
		desc = true
	case "asc", "ascending", "":
		// default
	default:
		return nil, fmt.Errorf("sort: unknown direction %q, want asc or desc", direction)
	}
	return &Sorter{field: field, descending: desc}, nil
}

// Apply buffers the line for later flushing; returns nil (line suppressed).
func (s *Sorter) Apply(line string) []string {
	s.buffer = append(s.buffer, line)
	return nil
}

// Flush sorts the buffered lines and returns them in order.
func (s *Sorter) Flush() []string {
	lines := s.buffer
	s.buffer = nil

	sort.SliceStable(lines, func(i, j int) bool {
		vi := extractSortValue(lines[i], s.field)
		vj := extractSortValue(lines[j], s.field)
		less := vi < vj
		if s.descending {
			return !less && vi != vj
		}
		return less
	})
	return lines
}

// Reset clears the internal buffer.
func (s *Sorter) Reset() {
	s.buffer = nil
}

func extractSortValue(line, field string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return ""
	}
	for k, v := range obj {
		if strings.EqualFold(k, field) {
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				return s
			}
			return string(v)
		}
	}
	return ""
}

// ParseSortFlag parses "field:asc" or "field:desc" (colon-separated).
func ParseSortFlag(s string) (*Sorter, error) {
	if s == "" {
		return nil, fmt.Errorf("sort: flag value must not be empty")
	}
	parts := strings.SplitN(s, ":", 2)
	dir := "asc"
	if len(parts) == 2 {
		dir = parts[1]
	}
	return NewSorter(parts[0], dir)
}
