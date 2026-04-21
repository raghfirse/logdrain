package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenameRule maps an old key to a new key.
type RenameRule struct {
	From string
	To   string
}

// Renamer applies field rename rules to JSON log lines.
type Renamer struct {
	rules []RenameRule
}

// NewRenamer creates a Renamer from a list of rules.
func NewRenamer(rules []RenameRule) *Renamer {
	return &Renamer{rules: rules}
}

// Apply renames fields in a JSON line according to the configured rules.
// Non-JSON lines are returned unchanged.
func (r *Renamer) Apply(line string) string {
	if len(r.rules) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, rule := range r.rules {
		key := findKeyCaseInsensitive(obj, rule.From)
		if key == "" {
			continue
		}
		obj[rule.To] = obj[key]
		delete(obj, key)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// ParseRenameFlag parses a slice of "old:new" strings into RenameRules.
func ParseRenameFlag(values []string) ([]RenameRule, error) {
	rules := make([]RenameRule, 0, len(values))
	for _, v := range values {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("rename: invalid rule %q, expected old:new", v)
		}
		rules = append(rules, RenameRule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}
