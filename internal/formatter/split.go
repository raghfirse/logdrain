package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Split reads a JSON field whose value is an array and emits one line per
// element. Non-JSON lines and lines where the target field is not an array
// are passed through unchanged.
type Split struct {
	field string
}

// NewSplit creates a Split processor targeting the given field name.
// An error is returned if field is empty.
func NewSplit(field string) (*Split, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("split: field name must not be empty")
	}
	return &Split{field: field}, nil
}

// Apply returns zero or more output lines derived from the input.
// If the named field contains a JSON array, one line is returned per element,
// each being the original object with the field replaced by the element value.
// Otherwise a single-element slice containing the original line is returned.
func (s *Split) Apply(line string) []string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return []string{line}
	}

	raw, ok := findRawKeyCaseInsensitive(obj, s.field)
	if !ok {
		return []string{line}
	}

	var elems []json.RawMessage
	if err := json.Unmarshal(raw, &elems); err != nil {
		// field exists but is not an array — pass through
		return []string{line}
	}

	out := make([]string, 0, len(elems))
	for _, elem := range elems {
		copy := make(map[string]json.RawMessage, len(obj))
		for k, v := range obj {
			copy[k] = v
		}
		// replace the field with the individual element
		for k := range copy {
			if strings.EqualFold(k, s.field) {
				copy[k] = elem
				break
			}
		}
		b, err := json.Marshal(copy)
		if err != nil {
			continue
		}
		out = append(out, string(b))
	}
	return out
}

// ParseSplitFlag parses the value of the --split flag (a field name).
func ParseSplitFlag(value string) (*Split, error) {
	return NewSplit(value)
}
