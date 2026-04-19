package formatter

import (
	"fmt"
	"strings"
)

// Wrapper wraps long field values at a given column width.
type Wrapper struct {
	width  int
	indent string
}

// NewWrapper creates a Wrapper that wraps values at the given width.
// A width of 0 disables wrapping.
func NewWrapper(width int, indent string) *Wrapper {
	return &Wrapper{width: width, indent: indent}
}

// Wrap applies line-wrapping to a string value, inserting newlines and
// the configured indent prefix at each break point.
func (w *Wrapper) Wrap(s string) string {
	if w.width <= 0 || len(s) <= w.width {
		return s
	}
	var sb strings.Builder
	for len(s) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("\n")
			sb.WriteString(w.indent)
		}
		if len(s) <= w.width {
			sb.WriteString(s)
			break
		}
		// Try to break on a space within the window.
		break_at := w.width
		if idx := strings.LastIndex(s[:w.width], " "); idx > 0 {
			break_at = idx + 1
		}
		sb.WriteString(s[:break_at])
		s = s[break_at:]
	}
	return sb.String()
}

// ParseWrapFlag parses a wrap width from a CLI flag string.
// Accepts a positive integer or "0"/"off" to disable.
func ParseWrapFlag(s string) (int, error) {
	if s == "" || s == "0" || s == "off" {
		return 0, nil
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil || n < 0 {
		return 0, fmt.Errorf("wrap: invalid width %q: must be a non-negative integer", s)
	}
	return n, nil
}
