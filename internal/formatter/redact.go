package formatter

import (
	"encoding/json"
	"strings"
)

const redactedValue = "[REDACTED]"

// Redactor masks specified fields in a JSON log line.
type Redactor struct {
	keys map[string]struct{}
}

// NewRedactor creates a Redactor that masks the given field keys (case-insensitive).
func NewRedactor(keys []string) *Redactor {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}
	return &Redactor{keys: m}
}

// Apply returns the line with matching fields replaced by [REDACTED].
// Non-JSON lines are returned unchanged.
func (r *Redactor) Apply(line string) string {
	if len(r.keys) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	modified := false
	for k := range obj {
		if _, ok := r.keys[strings.ToLower(k)]; ok {
			obj[k] = json.RawMessage(`"` + redactedValue + `"`)
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

// ParseRedactFlag parses a comma-separated list of field names to redact.
func ParseRedactFlag(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
