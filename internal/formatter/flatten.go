package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Flattener collapses nested JSON objects into dot-notation keys.
type Flattener struct {
	separator string
	prefix    string
}

// NewFlattener creates a Flattener with the given separator (e.g. ".").
func NewFlattener(separator string) *Flattener {
	if separator == "" {
		separator = "."
	}
	return &Flattener{separator: separator}
}

// Apply flattens nested keys in a JSON line. Non-JSON lines pass through.
func (f *Flattener) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	flat := make(map[string]interface{})
	f.flatten(obj, "", flat)
	b, err := json.Marshal(flat)
	if err != nil {
		return line
	}
	return string(b)
}

func (f *Flattener) flatten(obj map[string]interface{}, prefix string, out map[string]interface{}) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + f.separator + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			f.flatten(child, key, out)
		default:
			out[key] = v
		}
	}
}

// ParseFlattenFlag parses the --flatten flag value.
// Accepts a separator string or empty string to use the default ".".
func ParseFlattenFlag(s string) (*Flattener, error) {
	if strings.ContainsAny(s, " \t\n") {
		return nil, fmt.Errorf("flatten: separator must not contain whitespace")
	}
	return NewFlattener(s), nil
}
