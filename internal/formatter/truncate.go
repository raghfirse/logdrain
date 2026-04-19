package formatter

import "fmt"

// Truncator truncates field values that exceed a maximum length.
type Truncator struct {
	maxLen int
	suffix string
}

// NewTruncator creates a Truncator that clips string values to maxLen runes,
// appending suffix when truncation occurs. Pass 0 to disable truncation.
func NewTruncator(maxLen int, suffix string) *Truncator {
	if suffix == "" {
		suffix = "…"
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}
}

// Apply returns a copy of fields with string values truncated.
func (t *Truncator) Apply(fields map[string]any) map[string]any {
	if t.maxLen <= 0 {
		return fields
	}
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		if s, ok := v.(string); ok {
			runes := []rune(s)
			if len(runes) > t.maxLen {
				v = string(runes[:t.maxLen]) + t.suffix
			}
		}
		out[k] = v
	}
	return out
}

// ParseTruncateFlag parses a flag value like "120" or "120:..." into maxLen and suffix.
func ParseTruncateFlag(s string) (int, string, error) {
	if s == "" || s == "0" {
		return 0, "", nil
	}
	var maxLen int
	suffix := "…"
	n, err := fmt.Sscanf(s, "%d", &maxLen)
	if err != nil || n != 1 || maxLen < 1 {
		return 0, "", fmt.Errorf("truncate: invalid value %q, expected positive integer", s)
	}
	// Check for optional custom suffix after colon.
	for i, c := range s {
		if c == ':' {
			suffix = s[i+1:]
			break
		}
	}
	return maxLen, suffix, nil
}
