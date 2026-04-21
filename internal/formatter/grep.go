package formatter

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Grep filters log lines by matching a regex against one or more JSON fields.
// If no fields are specified, the entire raw line is matched.
type Grep struct {
	pattern *regexp.Regexp
	fields  []string
}

// NewGrep compiles pattern and returns a Grep filter scoped to the given fields.
// If fields is empty, the whole line is searched.
func NewGrep(pattern string, fields []string) (*Grep, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	norm := make([]string, len(fields))
	for i, f := range fields {
		norm[i] = strings.ToLower(f)
	}
	return &Grep{pattern: re, fields: norm}, nil
}

// Match reports whether line satisfies the grep filter.
func (g *Grep) Match(line string) bool {
	if len(g.fields) == 0 {
		return g.pattern.MatchString(line)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	for k, v := range obj {
		if !g.inFields(strings.ToLower(k)) {
			continue
		}
		if g.pattern.MatchString(stringify(v)) {
			return true
		}
	}
	return false
}

func (g *Grep) inFields(key string) bool {
	for _, f := range g.fields {
		if f == key {
			return true
		}
	}
	return false
}

func stringify(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case nil:
		return ""
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}
