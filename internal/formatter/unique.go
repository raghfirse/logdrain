package formatter

import (
	"encoding/json"
	"strings"
)

// Unique filters out lines where the value of a given field has already been seen.
// It is useful for suppressing repeated log entries that share a common identifier.
type Unique struct {
	field string
	seen  map[string]struct{}
}

// NewUnique creates a Unique filter for the given field name.
// Returns an error if field is empty.
func NewUnique(field string) (*Unique, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("unique: field name must not be empty")
	}
	return &Unique{
		field: strings.ToLower(field),
		seen:  make(map[string]struct{}),
	}, nil
}

// IsDuplicate returns true if the field value in line has been seen before.
// Non-JSON lines are never considered duplicates and always pass through.
func (u *Unique) IsDuplicate(line string) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}

	val := findRawValueCaseInsensitive(obj, u.field)
	if val == "" {
		return false
	}

	if _, exists := u.seen[val]; exists {
		return true
	}
	u.seen[val] = struct{}{}
	return false
}

// Reset clears all previously seen values.
func (u *Unique) Reset() {
	u.seen = make(map[string]struct{})
}

// ParseUniqueFlag returns the field name from the flag value, trimmed.
func ParseUniqueFlag(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("unique: field name must not be empty")
	}
	return s, nil
}

// findRawValueCaseInsensitive finds a key case-insensitively and returns its
// compact JSON representation as a string key for deduplication.
func findRawValueCaseInsensitive(obj map[string]json.RawMessage, key string) string {
	for k, v := range obj {
		if strings.ToLower(k) == key {
			return string(v)
		}
	}
	return ""
}
