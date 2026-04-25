package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// After suppresses log lines whose timestamp is after a given cutoff time.
// Lines without a recognizable timestamp are passed through unchanged.
type After struct {
	cutoff time.Time
}

// NewAfter creates an After filter from a cutoff time.
// The cutoff may be an RFC3339 string or a negative duration like "-1h".
func NewAfter(cutoff time.Time) (*After, error) {
	if cutoff.IsZero() {
		return nil, fmt.Errorf("after: cutoff time must not be zero")
	}
	return &After{cutoff: cutoff}, nil
}

// ParseAfterFlag parses a flag value into an After filter.
// Accepts RFC3339 timestamps or relative durations (e.g. "-30m").
func ParseAfterFlag(s string) (*After, error) {
	if s == "" {
		return nil, fmt.Errorf("after: value must not be empty")
	}
	// Try relative duration first (e.g. "-1h", "30m")
	if strings.HasPrefix(s, "-") || (!strings.Contains(s, "T") && !strings.Contains(s, "-")) {
		d, err := time.ParseDuration(s)
		if err == nil {
			if d > 0 {
				return nil, fmt.Errorf("after: relative duration must be negative (e.g. -1h)")
			}
			return NewAfter(time.Now().Add(d))
		}
	}
	// Try RFC3339
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return NewAfter(t)
	}
	return nil, fmt.Errorf("after: cannot parse %q as RFC3339 or duration", s)
}

// Apply returns the line if its timestamp is before or equal to the cutoff,
// suppressing lines that are after the cutoff. Non-JSON lines pass through.
func (a *After) Apply(line string) (string, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true
	}
	t := extractAfterTimestamp(obj)
	if t.IsZero() {
		return line, true
	}
	if t.After(a.cutoff) {
		return "", false
	}
	return line, true
}

func extractAfterTimestamp(obj map[string]json.RawMessage) time.Time {
	keys := []string{"time", "ts", "timestamp", "@timestamp"}
	for _, k := range keys {
		for objKey, raw := range obj {
			if !strings.EqualFold(objKey, k) {
				continue
			}
			var s string
			if err := json.Unmarshal(raw, &s); err == nil {
				if t, err := time.Parse(time.RFC3339, s); err == nil {
					return t
				}
				if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
					return t
				}
			}
			var f float64
			if err := json.Unmarshal(raw, &f); err == nil {
				sec := int64(f)
				ns := int64((f - float64(sec)) * 1e9)
				return time.Unix(sec, ns)
			}
			var n json.Number
			if err := json.Unmarshal(raw, &n); err == nil {
				if i, err := strconv.ParseInt(n.String(), 10, 64); err == nil {
					return time.Unix(i, 0)
				}
			}
		}
	}
	return time.Time{}
}
