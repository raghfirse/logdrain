package formatter

import (
	"strings"
	"time"
)

// Common timestamp field names to probe.
var timestampKeys = []string{"time", "timestamp", "ts", "@timestamp"}

// knownFormats are tried in order when parsing timestamp strings.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	time.UnixDate,
}

// FormatTimestamp extracts a timestamp from a log entry map and returns a
// human-readable string. Returns empty string if no timestamp is found.
func FormatTimestamp(entry map[string]any, short bool) string {
	for _, key := range timestampKeys {
		val, ok := entry[key]
		if !ok {
			continue
		}
		if t, ok := parseTimestampValue(val); ok {
			if short {
				return t.Format("15:04:05")
			}
			return t.Format("2006-01-02 15:04:05.000")
		}
	}
	return ""
}

func parseTimestampValue(val any) (time.Time, bool) {
	switch v := val.(type) {
	case string:
		return parseTimestampString(v)
	case float64:
		// Unix epoch seconds (with possible fractional)
		sec := int64(v)
		nsec := int64((v - float64(sec)) * 1e9)
		return time.Unix(sec, nsec).UTC(), true
	}
	return time.Time{}, false
}

func parseTimestampString(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, fmt := range knownFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}
