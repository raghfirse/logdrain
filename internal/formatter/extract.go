package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extractor promotes a nested JSON field to the top level.
// For example, given key "meta.request_id", it reads
// obj["meta"]["request_id"] and writes it as obj["request_id"].
type Extractor struct {
	rules []extractRule
}

type extractRule struct {
	path   []string // dot-separated path segments
	alias  string   // top-level key to write (last segment by default)
}

// NewExtractor parses one or more extract rules of the form
// "a.b.c" or "a.b.c->alias" and returns an Extractor.
func NewExtractor(specs []string) (*Extractor, error) {
	if len(specs) == 0 {
		return &Extractor{}, nil
	}
	rules := make([]extractRule, 0, len(specs))
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		alias := ""
		if idx := strings.Index(spec, "->"); idx >= 0 {
			alias = strings.TrimSpace(spec[idx+2:])
			spec = strings.TrimSpace(spec[:idx])
		}
		parts := strings.Split(spec, ".")
		for _, p := range parts {
			if p == "" {
				return nil, fmt.Errorf("extract: invalid path %q", spec)
			}
		}
		if alias == "" {
			alias = parts[len(parts)-1]
		}
		rules = append(rules, extractRule{path: parts, alias: alias})
	}
	return &Extractor{rules: rules}, nil
}

// Apply reads line, extracts nested fields according to its rules,
// and returns the modified JSON line. Non-JSON lines are returned as-is.
func (e *Extractor) Apply(line string) string {
	if len(e.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range e.rules {
		val, ok := nestedGet(obj, r.path)
		if !ok {
			continue
		}
		obj[r.alias] = val
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// nestedGet walks obj following path segments (case-insensitive) and
// returns the value at the leaf, or (nil, false) if not found.
func nestedGet(obj map[string]interface{}, path []string) (interface{}, bool) {
	var cur interface{} = obj
	for _, seg := range path {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, false
		}
		key := findKeyCaseInsensitive(m, seg)
		if key == "" {
			return nil, false
		}
		cur = m[key]
	}
	return cur, true
}
