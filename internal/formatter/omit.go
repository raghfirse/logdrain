package formatter

import (
	"encoding/json"
	"strings"
)

// Omit removes specified keys from JSON log lines.
type Omit struct {
	keys []string
}

// NewOmit creates an Omit that removes the given keys (case-insensitive) from
// each JSON object. Returns an error if no keys are provided.
func NewOmit(keys []string) (*Omit, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("omit: at least one key is required")
	}
	norm := make([]string, len(keys))
	for i, k := range keys {
		norm[i] = strings.ToLower(k)
	}
	return &Omit{keys: norm}, nil
}

// Apply removes the configured keys from the JSON line. Non-JSON lines are
// passed through unchanged.
func (o *Omit) Apply(line string) (string, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true
	}
	for rawKey := range obj {
		for _, k := range o.keys {
			if strings.ToLower(rawKey) == k {
				delete(obj, rawKey)
				break
			}
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line, true
	}
	return string(out), true
}

// ParseOmitFlag parses a comma-separated list of keys to omit.
func ParseOmitFlag(s string) ([]string, error) {
	if strings.TrimSpace(s) == "" {
		return nil, fmt.Errorf("omit: flag value must not be empty")
	}
	parts := strings.Split(s, ",")
	keys := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			keys = append(keys, p)
		}
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("omit: no valid keys found in %q", s)
	}
	return keys, nil
}
