package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Clamp constrains a numeric field's value to a [min, max] range.
type Clamp struct {
	key string
	min float64
	max float64
}

// NewClamp creates a Clamp transformer for the given field and bounds.
// Returns an error if min > max or key is empty.
func NewClamp(key string, min, max float64) (*Clamp, error) {
	if key == "" {
		return nil, fmt.Errorf("clamp: key must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("clamp: min %.g is greater than max %.g", min, max)
	}
	return &Clamp{key: key, min: min, max: max}, nil
}

// Apply clamps the target field in a JSON log line.
// Non-JSON lines and lines missing the field are passed through unchanged.
func (c *Clamp) Apply(line string) (string, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true
	}

	raw, ok := findRawKeyCaseInsensitive(obj, c.key)
	if !ok {
		return line, true
	}

	v, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		return line, true
	}

	if v < c.min {
		v = c.min
	} else if v > c.max {
		v = c.max
	}

	obj[c.key] = json.RawMessage(strconv.FormatFloat(v, 'f', -1, 64))

	out, err := json.Marshal(obj)
	if err != nil {
		return line, true
	}
	return string(out), true
}

// ParseClampFlag parses a clamp spec of the form "field,min,max".
func ParseClampFlag(s string) (*Clamp, error) {
	parts := splitN(s, ",", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("clamp: expected format field,min,max — got %q", s)
	}
	min, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("clamp: invalid min %q: %w", parts[1], err)
	}
	max, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("clamp: invalid max %q: %w", parts[2], err)
	}
	return NewClamp(parts[0], min, max)
}

func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx < 0 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	return append(parts, s)
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
