package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Script applies a simple key=expr rewrite to a JSON log line.
// Each rule is of the form "key=gotemplate" where the template
// has access to all top-level JSON fields.
type Script struct {
	rules []scriptRule
}

type scriptRule struct {
	key  string
	tmpl *Template
}

// NewScript parses a slice of "key=template" rule strings and returns a Script.
func NewScript(rules []string) (*Script, error) {
	parsed := make([]scriptRule, 0, len(rules))
	for _, r := range rules {
		idx := strings.IndexByte(r, '=')
		if idx < 1 {
			return nil, fmt.Errorf("script: invalid rule %q: expected key=template", r)
		}
		key := strings.TrimSpace(r[:idx])
		expr := strings.TrimSpace(r[idx+1:])
		if key == "" {
			return nil, fmt.Errorf("script: empty key in rule %q", r)
		}
		tmpl, err := NewTemplate(expr)
		if err != nil {
			return nil, fmt.Errorf("script: rule %q: %w", r, err)
		}
		parsed = append(parsed, scriptRule{key: key, tmpl: tmpl})
	}
	return &Script{rules: parsed}, nil
}

// Apply evaluates each rule against the JSON line and returns the rewritten line.
// Non-JSON lines are returned unchanged.
func (s *Script) Apply(line string) string {
	if len(s.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, rule := range s.rules {
		val, err := rule.tmpl.Render(line)
		if err != nil {
			continue
		}
		obj[rule.key] = val
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
