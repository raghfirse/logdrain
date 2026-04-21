package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Merger merges fields from a static JSON object into each log line.
type Merger struct {
	fields map[string]json.RawMessage
	overwrite bool
}

// NewMerge creates a Merger from a JSON object string.
// If overwrite is true, merged fields replace existing keys.
func NewMerge(jsonObj string, overwrite bool) (*Merger, error) {
	if jsonObj == "" {
		return &Merger{fields: map[string]json.RawMessage{}, overwrite: overwrite}, nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(jsonObj), &fields); err != nil {
		return nil, fmt.Errorf("merge: invalid JSON object: %w", err)
	}
	return &Merger{fields: fields, overwrite: overwrite}, nil
}

// Apply merges the static fields into the given log line.
func (m *Merger) Apply(line string) string {
	if len(m.fields) == 0 {
		return line
	}
	trimmed := strings.TrimSpace(line)
	var record map[string]json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &record); err != nil {
		return line
	}
	for k, v := range m.fields {
		if _, exists := record[k]; exists && !m.overwrite {
			continue
		}
		record[k] = v
	}
	out, err := json.Marshal(record)
	if err != nil {
		return line
	}
	return string(out)
}

// ParseMergeFlag parses a --merge flag value of the form '{"key":"val"}'.
func ParseMergeFlag(s string) (*Merger, error) {
	return NewMerge(s, false)
}
