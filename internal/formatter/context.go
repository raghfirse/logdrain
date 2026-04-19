package formatter

import (
	"encoding/json"
	"strings"
)

// ContextFields extracts a subset of fields from a JSON log line
// to display as contextual metadata alongside the main message.
type ContextFields struct {
	keys map[string]struct{}
}

// NewContextFields creates a ContextFields extractor from a list of keys.
func NewContextFields(keys []string) *ContextFields {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}
	return &ContextFields{keys: m}
}

// Extract returns a map of matching key/value pairs from raw JSON.
// Returns nil if no keys match or input is not valid JSON.
func (c *ContextFields) Extract(line []byte) map[string]any {
	if len(c.keys) == 0 {
		return nil
	}
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil
	}
	out := make(map[string]any)
	for k, v := range obj {
		if _, ok := c.keys[strings.ToLower(k)]; ok {
			out[k] = v
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// ParseContextFlag parses a comma-separated list of field names.
func ParseContextFlag(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
