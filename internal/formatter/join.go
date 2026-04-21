package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JoinRule defines how to combine multiple fields into a single output field.
type JoinRule struct {
	Keys      []string
	OutputKey string
	Sep       string
}

// Join combines multiple JSON fields into a single field using a separator.
type Join struct {
	rules []JoinRule
}

// NewJoin creates a Join transformer from a slice of rules.
func NewJoin(rules []JoinRule) *Join {
	return &Join{rules: rules}
}

// Apply processes a log line, merging fields per each rule.
func (j *Join) Apply(line string) string {
	if len(j.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, rule := range j.rules {
		parts := make([]string, 0, len(rule.Keys))
		for _, k := range rule.Keys {
			if v, ok := findKeyCaseInsensitive(obj, k); ok {
				parts = append(parts, fmt.Sprintf("%v", v))
			}
		}
		obj[rule.OutputKey] = strings.Join(parts, rule.Sep)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// ParseJoinFlag parses a join flag string of the form:
//   "key1+key2->output" or "key1+key2:sep->output"
func ParseJoinFlag(s string) (JoinRule, error) {
	const arrow = "->"
	idx := strings.Index(s, arrow)
	if idx < 0 {
		return JoinRule{}, fmt.Errorf("join: missing '->': %q", s)
	}
	lhs := s[:idx]
	output := strings.TrimSpace(s[idx+len(arrow):])
	if output == "" {
		return JoinRule{}, fmt.Errorf("join: empty output key in %q", s)
	}
	sep := " "
	if ci := strings.LastIndex(lhs, ":"); ci >= 0 {
		sep = lhs[ci+1:]
		lhs = lhs[:ci]
	}
	rawKeys := strings.Split(lhs, "+")
	keys := make([]string, 0, len(rawKeys))
	for _, k := range rawKeys {
		k = strings.TrimSpace(k)
		if k == "" {
			return JoinRule{}, fmt.Errorf("join: empty key in %q", s)
		}
		keys = append(keys, k)
	}
	if len(keys) < 2 {
		return JoinRule{}, fmt.Errorf("join: need at least two keys in %q", s)
	}
	return JoinRule{Keys: keys, OutputKey: output, Sep: sep}, nil
}
