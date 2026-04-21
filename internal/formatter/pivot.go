package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Pivot restructures a JSON log line by promoting a nested field's value as the
// top-level key, grouping the remaining fields under it. This is useful when
// logs contain a discriminator field (e.g. "event" or "type") and you want to
// reshape the output around that value.
//
// Example with key="event":
//
//	input:  {"event":"user.login","user":"alice","ip":"1.2.3.4"}
//	output: {"user.login":{"user":"alice","ip":"1.2.3.4"}}
type Pivot struct {
	key    string
	exclude bool
}

// NewPivot creates a Pivot that restructures lines around the given key.
// If exclude is true, the pivot key itself is omitted from the nested object.
func NewPivot(key string, exclude bool) (*Pivot, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("pivot: key must not be empty")
	}
	return &Pivot{key: key, exclude: exclude}, nil
}

// Apply pivots a single log line. Non-JSON lines are returned unchanged.
// If the pivot key is missing or its value is not a string, the line is
// returned unchanged.
func (p *Pivot) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	// Find the pivot key (case-insensitive).
	actualKey, raw, ok := findRawKeyCaseInsensitive(obj, p.key)
	if !ok {
		return line
	}

	// The pivot value must be a plain string.
	var pivotValue string
	if err := json.Unmarshal(raw, &pivotValue); err != nil {
		return line
	}

	// Build the inner object, optionally dropping the pivot key.
	inner := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		if p.exclude && strings.EqualFold(k, actualKey) {
			continue
		}
		inner[k] = v
	}

	innerBytes, err := json.Marshal(inner)
	if err != nil {
		return line
	}

	outer := map[string]json.RawMessage{
		pivotValue: innerBytes,
	}
	out, err := json.Marshal(outer)
	if err != nil {
		return line
	}
	return string(out)
}

// findRawKeyCaseInsensitive looks up a key in a raw JSON object map using
// case-insensitive matching and returns the actual key name, raw value, and
// whether it was found.
func findRawKeyCaseInsensitive(obj map[string]json.RawMessage, key string) (string, json.RawMessage, bool) {
	// Exact match first.
	if v, ok := obj[key]; ok {
		return key, v, true
	}
	// Fall back to case-insensitive scan.
	for k, v := range obj {
		if strings.EqualFold(k, key) {
			return k, v, true
		}
	}
	return "", nil, false
}

// ParsePivotFlag parses a pivot flag value of the form "key" or "key,exclude".
// The optional ",exclude" suffix causes the pivot key to be dropped from the
// nested object.
func ParsePivotFlag(s string) (*Pivot, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("pivot: flag value must not be empty")
	}
	parts := strings.SplitN(s, ",", 2)
	key := strings.TrimSpace(parts[0])
	exclude := false
	if len(parts) == 2 {
		mod := strings.TrimSpace(strings.ToLower(parts[1]))
		if mod != "exclude" {
			return nil, fmt.Errorf("pivot: unknown modifier %q (expected \"exclude\")", parts[1])
		}
		exclude = true
	}
	return NewPivot(key, exclude)
}
