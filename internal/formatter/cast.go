package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// CastRule defines a field and the target type to cast it to.
type CastRule struct {
	Key  string
	Type string // "string", "int", "float", "bool"
}

// Cast applies type coercion rules to JSON log lines.
type Cast struct {
	rules []CastRule
}

// NewCast creates a Cast transformer from a list of rules.
func NewCast(rules []CastRule) (*Cast, error) {
	for _, r := range rules {
		switch r.Type {
		case "string", "int", "float", "bool":
			// valid
		default:
			return nil, fmt.Errorf("cast: unknown type %q for key %q", r.Type, r.Key)
		}
		if r.Key == "" {
			return nil, fmt.Errorf("cast: key must not be empty")
		}
	}
	return &Cast{rules: rules}, nil
}

// Apply coerces fields in a JSON log line according to the configured rules.
// Non-JSON lines are returned unchanged.
func (c *Cast) Apply(line string) string {
	if len(c.rules) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, rule := range c.rules {
		key := findKeyCaseInsensitive(obj, rule.Key)
		if key == "" {
			continue
		}
		raw := strings.TrimSpace(string(obj[key]))
		// Strip surrounding quotes for string values
		unquoted := raw
		if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
			unquoted = raw[1 : len(raw)-1]
		}
		var newVal json.RawMessage
		switch rule.Type {
		case "string":
			b, _ := json.Marshal(unquoted)
			newVal = b
		case "int":
			f, err := strconv.ParseFloat(unquoted, 64)
			if err != nil {
				continue
			}
			newVal = json.RawMessage(strconv.FormatInt(int64(f), 10))
		case "float":
			f, err := strconv.ParseFloat(unquoted, 64)
			if err != nil {
				continue
			}
			b, _ := json.Marshal(f)
			newVal = b
		case "bool":
			v, err := strconv.ParseBool(unquoted)
			if err != nil {
				continue
			}
			b, _ := json.Marshal(v)
			newVal = b
		}
		obj[key] = newVal
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// ParseCastFlag parses a cast flag value in the form "key:type".
func ParseCastFlag(s string) (CastRule, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return CastRule{}, fmt.Errorf("cast: invalid format %q, expected key:type", s)
	}
	return CastRule{Key: parts[0], Type: parts[1]}, nil
}
