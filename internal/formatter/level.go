package formatter

import (
	"fmt"
	"strings"
)

// LevelOrder defines severity order for known log levels.
var LevelOrder = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
}

// NormalizeLevel returns a lowercase, canonical level string.
func NormalizeLevel(level string) string {
	l := strings.ToLower(strings.TrimSpace(level))
	switch l {
	case "warning":
		return "warn"
	case "err":
		return "error"
	case "critical", "panic":
		return "fatal"
	}
	return l
}

// ParseLevelFilter parses a minimum level string and returns a filter func.
// The returned func returns true if the given level meets the minimum.
func ParseLevelFilter(minLevel string) (func(string) bool, error) {
	norm := NormalizeLevel(minLevel)
	minVal, ok := LevelOrder[norm]
	if !ok {
		return nil, fmt.Errorf("unknown log level: %q", minLevel)
	}
	return func(level string) bool {
		v, ok := LevelOrder[NormalizeLevel(level)]
		if !ok {
			return true
		}
		return v >= minVal
	}, nil
}
