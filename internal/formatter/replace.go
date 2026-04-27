package formatter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ReplaceRule describes a single field value replacement using a regex.
type ReplaceRule struct {
	Key     string
	Pattern *regexp.Regexp
	With    string
}

// Replacer applies regex-based value replacements to JSON log fields.
type Replacer struct {
	rules []ReplaceRule
}

// NewReplace constructs a Replacer from a slice of ReplaceRule values.
func NewReplace(rules []ReplaceRule) (*Replacer, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("replace: at least one rule is required")
	}
	return &Replacer{rules: rules}, nil
}

// Apply performs replacements on matching fields in the JSON line.
// Non-JSON lines are passed through unchanged.
func (r *Replacer) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, rule := range r.rules {
		key := findKeyCaseInsensitive(obj, rule.Key)
		if key == "" {
			continue
		}
		var s string
		if err := json.Unmarshal(obj[key], &s); err != nil {
			continue
		}
		replaced := rule.Pattern.ReplaceAllString(s, rule.With)
		encoded, err := json.Marshal(replaced)
		if err != nil {
			continue
		}
		obj[key] = json.RawMessage(encoded)
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// ParseReplaceFlag parses a flag value of the form "key/pattern/replacement".
func ParseReplaceFlag(s string) (ReplaceRule, error) {
	parts := strings.SplitN(s, "/", 3)
	if len(parts) != 3 {
		return ReplaceRule{}, fmt.Errorf("replace: expected key/pattern/replacement, got %q", s)
	}
	key := strings.TrimSpace(parts[0])
	if key == "" {
		return ReplaceRule{}, fmt.Errorf("replace: key must not be empty")
	}
	re, err := regexp.Compile(parts[1])
	if err != nil {
		return ReplaceRule{}, fmt.Errorf("replace: invalid pattern %q: %w", parts[1], err)
	}
	return ReplaceRule{Key: key, Pattern: re, With: parts[2]}, nil
}
