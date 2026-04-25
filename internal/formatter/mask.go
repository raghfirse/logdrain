package formatter

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Masker replaces values of matching fields with a pattern-based mask,
// useful for partially obscuring sensitive data (e.g. showing last 4 digits).
type Masker struct {
	rules []maskRule
}

type maskRule struct {
	key     string
	pattern *regexp.Regexp
	mask    string
}

// NewMask creates a Masker from a slice of rule specs.
// Each spec has the form "field:pattern:mask", e.g. "card:[0-9]{12}(\\d{4}):****$1".
func NewMask(specs []string) (*Masker, error) {
	rules := make([]maskRule, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("mask: invalid spec %q: expected field:pattern:mask", spec)
		}
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return nil, fmt.Errorf("mask: invalid pattern %q: %w", parts[1], err)
		}
		rules = append(rules, maskRule{key: strings.ToLower(parts[0]), pattern: re, mask: parts[2]})
	}
	return &Masker{rules: rules}, nil
}

// Apply rewrites matching string fields using the configured regex replacement.
func (m *Masker) Apply(line string) string {
	if len(m.rules) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	modified := false
	for k, raw := range obj {
		for _, rule := range m.rules {
			if strings.ToLower(k) != rule.key {
				continue
			}
			var s string
			if err := json.Unmarshal(raw, &s); err != nil {
				continue
			}
			replaced := rule.pattern.ReplaceAllString(s, rule.mask)
			if replaced == s {
				continue
			}
			b, _ := json.Marshal(replaced)
			obj[k] = b
			modified = true
		}
	}
	if !modified {
		return line
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// ParseMaskFlag converts a slice of raw flag values into a Masker.
func ParseMaskFlag(vals []string) (*Masker, error) {
	return NewMask(vals)
}
