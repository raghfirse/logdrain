package formatter

import (
	"encoding/json"
	"fmt"
	"time"
)

// Since filters log lines, only passing through entries whose timestamp
// field is at or after the given cutoff time.
type Since struct {
	cutoff time.Time
	field  string
}

// NewSince creates a Since filter that passes lines with a timestamp >= cutoff.
// field is the JSON key to inspect; if empty, common aliases are tried.
func NewSince(cutoff time.Time, field string) *Since {
	if field == "" {
		field = "time"
	}
	return &Since{cutoff: cutoff, field: field}
}

// Apply returns the line if its timestamp is at or after the cutoff,
// an empty string if it should be suppressed, or the original line
// if it cannot be parsed.
func (s *Since) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	raw, ok := findRawKeyCaseInsensitive(obj, s.field)
	if !ok {
		// Try common aliases.
		for _, alias := range []string{"ts", "timestamp", "@timestamp", "t"} {
			if r, found := findRawKeyCaseInsensitive(obj, alias); found {
				raw = r
				ok = true
				break
			}
		}
	}
	if !ok {
		return line
	}

	t, err := parseTimestampValue(raw)
	if err != nil {
		return line
	}

	if t.Before(s.cutoff) {
		return ""
	}
	return line
}

// ParseSinceFlag parses a duration string like "5m", "1h", "30s" and
// returns a cutoff time of now minus that duration.
func ParseSinceFlag(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid since duration %q: %w", s, err)
	}
	if d < 0 {
		return time.Time{}, fmt.Errorf("since duration must be positive, got %q", s)
	}
	return time.Now().Add(-d), nil
}
