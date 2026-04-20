package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Transform applies a set of key rename/copy operations to a JSON log line.
// Each rule is in the form "old:new" — the value at "old" is copied to "new",
// and the original key is removed.
type Transform struct {
	rules []transformRule
}

type transformRule struct {
	from string
	to   string
}

// NewTransform creates a Transform from a slice of "from:to" rule strings.
func NewTransform(rules []string) (*Transform, error) {
	parsed := make([]transformRule, 0, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid transform rule %q: expected \"from:to\"", r)
		}
		parsed = append(parsed, transformRule{from: parts[0], to: parts[1]})
	}
	return &Transform{rules: parsed}, nil
}

// Apply renames fields in a JSON line according to the configured rules.
// Non-JSON lines are returned unchanged. If no rules are configured the
// original line is returned as-is.
func (t *Transform) Apply(line string) string {
	if len(t.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, rule := range t.rules {
		key := findKeyCaseInsensitive(obj, rule.from)
		if key == "" {
			continue
		}
		obj[rule.to] = obj[key]
		delete(obj, key)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// findKeyCaseInsensitive returns the actual map key that matches name
// case-insensitively, or empty string if not found.
func findKeyCaseInsensitive(obj map[string]interface{}, name string) string {
	lower := strings.ToLower(name)
	for k := range obj {
		if strings.ToLower(k) == lower {
			return k
		}
	}
	return ""
}
