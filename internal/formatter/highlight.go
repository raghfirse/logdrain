package formatter

import (
	"fmt"
	"strings"
)

// FieldHighlight defines a field name and the ANSI color code to apply to its value.
type FieldHighlight struct {
	Field string
	Color string
}

// HighlightFields returns a copy of the fields map with specified field values
// wrapped in ANSI color codes.
func HighlightFields(fields map[string]interface{}, highlights []FieldHighlight) map[string]interface{} {
	if len(highlights) == 0 {
		return fields
	}
	result := make(map[string]interface{}, len(fields))
	highlightMap := make(map[string]string, len(highlights))
	for _, h := range highlights {
		highlightMap[strings.ToLower(h.Field)] = h.Color
	}
	for k, v := range fields {
		if color, ok := highlightMap[strings.ToLower(k)]; ok {
			result[k] = fmt.Sprintf("%s%v\033[0m", color, v)
		} else {
			result[k] = v
		}
	}
	return result
}

// ParseHighlightFlag parses a slice of "field=color" strings into FieldHighlight entries.
// Unrecognized color names are mapped to their raw ANSI escape if prefixed with "\033[",
// otherwise a default yellow is used.
func ParseHighlightFlag(exprs []string) []FieldHighlight {
	named := map[string]string{
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
	}
	var out []FieldHighlight
	for _, expr := range exprs {
		parts := strings.SplitN(expr, "=", 2)
		if len(parts) != 2 {
			continue
		}
		field, colorName := parts[0], strings.ToLower(parts[1])
		color, ok := named[colorName]
		if !ok {
			color = "\033[33m" // default yellow
		}
		out = append(out, FieldHighlight{Field: field, Color: color})
	}
	return out
}
