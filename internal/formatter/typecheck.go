package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TypeChecker filters or flags log lines where a field's JSON type does not
// match the expected type. Supported type names: string, number, bool, array,
// object, null.
type TypeChecker struct {
	key      string
	wantType string
	invert   bool
}

// NewTypeChecker creates a TypeChecker that passes lines where the named field
// has the given JSON type. If invert is true, lines that do NOT match are
// passed instead.
func NewTypeChecker(key, typeName string, invert bool) (*TypeChecker, error) {
	if key == "" {
		return nil, fmt.Errorf("typecheck: key must not be empty")
	}
	norm := strings.ToLower(typeName)
	switch norm {
	case "string", "number", "bool", "array", "object", "null":
	default:
		return nil, fmt.Errorf("typecheck: unknown type %q; want string|number|bool|array|object|null", typeName)
	}
	return &TypeChecker{key: key, wantType: norm, invert: invert}, nil
}

// Apply returns the line unchanged if the field type matches (or does not match
// when inverted). Returns an empty string to suppress the line.
func (t *TypeChecker) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	val, ok := findRawValueCaseInsensitive(obj, t.key)
	if !ok {
		if t.invert {
			return line
		}
		return ""
	}
	actual := jsonTypeName(val)
	matches := actual == t.wantType
	if t.invert {
		matches = !matches
	}
	if matches {
		return line
	}
	return ""
}

// jsonTypeName returns the JSON type name for a raw JSON value.
func jsonTypeName(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "null"
	}
	switch raw[0] {
	case '"':
		return "string"
	case '{':
		return "object"
	case '[':
		return "array"
	case 't', 'f':
		return "bool"
	case 'n':
		return "null"
	default:
		return "number"
	}
}

// ParseTypeCheckFlag parses a flag value of the form "key:type" or
// "!key:type" (inverted). Returns a configured TypeChecker.
func ParseTypeCheckFlag(s string) (*TypeChecker, error) {
	invert := false
	if strings.HasPrefix(s, "!") {
		invert = true
		s = s[1:]
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("typecheck: expected key:type, got %q", s)
	}
	return NewTypeChecker(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), invert)
}
