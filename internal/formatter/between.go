package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Between suppresses lines until a start pattern is matched, then passes lines
// through until an end pattern is matched (inclusive of both boundary lines).
type Between struct {
	startPattern string
	endPattern   string
	field        string
	active       bool
}

// NewBetween creates a Between filter. startPattern and endPattern are
// substring matches applied to the resolved field value (or whole line).
func NewBetween(startPattern, endPattern, field string) (*Between, error) {
	if startPattern == "" {
		return nil, fmt.Errorf("between: start pattern must not be empty")
	}
	if endPattern == "" {
		return nil, fmt.Errorf("between: end pattern must not be empty")
	}
	return &Between{
		startPattern: startPattern,
		endPattern:   endPattern,
		field:        field,
	}, nil
}

// Apply returns the line if it falls within the active window, otherwise "".
// Flush returns nothing; Between is stateful but produces no deferred output.
func (b *Between) Apply(line string) string {
	value := b.resolve(line)

	if !b.active {
		if strings.Contains(value, b.startPattern) {
			b.active = true
			return line
		}
		return ""
	}

	// Already active — pass line through, then check for end.
	if strings.Contains(value, b.endPattern) {
		b.active = false
	}
	return line
}

func (b *Between) Flush() []string { return nil }

func (b *Between) Reset() { b.active = false }

func (b *Between) resolve(line string) string {
	if b.field == "" {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for k, v := range obj {
		if strings.EqualFold(k, b.field) {
			var s string
			if err := json.Unmarshal(v, &s); err == nil {
				return s
			}
			return string(v)
		}
	}
	return ""
}

// ParseBetweenFlag parses a flag value of the form "start:end" or
// "start:end:field".
func ParseBetweenFlag(s string) (*Between, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("between: expected start:end[:field], got %q", s)
	}
	field := ""
	if len(parts) == 3 {
		field = parts[2]
	}
	return NewBetween(parts[0], parts[1], field)
}
