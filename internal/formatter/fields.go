package formatter

import "strings"

// FieldFilter controls which JSON fields are included in pretty output.
type FieldFilter struct {
	exclude map[string]struct{}
	include map[string]struct{}
}

// ParseFieldsFlag parses a comma-separated list of field names.
func ParseFieldsFlag(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, strings.ToLower(p))
		}
	}
	return out
}

// NewFieldFilter builds a FieldFilter from include/exclude lists.
// Exclude takes precedence over include.
func NewFieldFilter(include, exclude []string) *FieldFilter {
	ff := &FieldFilter{
		exclude: make(map[string]struct{}),
		include: make(map[string]struct{}),
	}
	for _, f := range include {
		ff.include[strings.ToLower(f)] = struct{}{}
	}
	for _, f := range exclude {
		ff.exclude[strings.ToLower(f)] = struct{}{}
	}
	return ff
}

// Allow reports whether the given field key should be shown.
func (ff *FieldFilter) Allow(key string) bool {
	lk := strings.ToLower(key)
	if _, excluded := ff.exclude[lk]; excluded {
		return false
	}
	if len(ff.include) == 0 {
		return true
	}
	_, included := ff.include[lk]
	return included
}

// IsEmpty returns true when no include/exclude rules are defined.
func (ff *FieldFilter) IsEmpty() bool {
	return len(ff.include) == 0 && len(ff.exclude) == 0
}
