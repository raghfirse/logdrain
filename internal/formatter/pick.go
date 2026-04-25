package formatter

import (
	"encoding/json"
	"strings"
)

// Pick retains only the specified top-level keys from each JSON log line.
// Non-JSON lines are passed through unchanged.
type Pick struct {
	keys []string
}

// NewPick creates a Pick transformer from a list of field names.
// Keys are matched case-insensitively.
func NewPick(keys []string) (*Pick, error) {
	if len(keys) == 0 {
		return &Pick{}, nil
	}
	normalized := make([]string, len(keys))
	for i, k := range keys {
		normalized[i] = strings.ToLower(strings.TrimSpace(k))
	}
	return &Pick{keys: normalized}, nil
}

// Apply filters the JSON object to include only the picked keys.
func (p *Pick) Apply(line string) (string, bool) {
	if len(p.keys) == 0 {
		return line, true
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true
	}
	out := make(map[string]json.RawMessage, len(p.keys))
	for rawKey, val := range obj {
		for _, want := range p.keys {
			if strings.ToLower(rawKey) == want {
				out[rawKey] = val
				break
			}
		}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return line, true
	}
	return string(b), true
}

// ParsePickFlag parses a comma-separated list of field names.
func ParsePickFlag(s string) ([]string, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result, nil
}
