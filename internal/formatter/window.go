package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Window groups log lines into fixed time buckets and emits a summary line
// when the bucket closes. Useful for rate-of-occurrence analysis.
type Window struct {
	field    string
	duration time.Duration
	buckets  map[string]int
	bucketAt time.Time
	now      func() time.Time
}

// NewWindow creates a Window that groups by the given field over the given duration.
func NewWindow(field string, d time.Duration) (*Window, error) {
	if field == "" {
		return nil, fmt.Errorf("window: field must not be empty")
	}
	if d <= 0 {
		return nil, fmt.Errorf("window: duration must be positive")
	}
	return &Window{
		field:    field,
		duration: d,
		buckets:  make(map[string]int),
		now:      time.Now,
	}, nil
}

// Apply records the line in the current bucket. If the bucket window has
// expired it flushes and returns a summary JSON line; otherwise returns "".
func (w *Window) Apply(line string) string {
	now := w.now()
	if w.bucketAt.IsZero() {
		w.bucketAt = now
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		key := ""
		for k, v := range obj {
			if strings.EqualFold(k, w.field) {
				key = fmt.Sprintf("%v", v)
				break
			}
		}
		w.buckets[key]++
	}

	if now.Sub(w.bucketAt) >= w.duration {
		return w.flush()
	}
	return ""
}

// Flush forces emission of the current bucket summary and resets state.
func (w *Window) Flush() string {
	return w.flush()
}

func (w *Window) flush() string {
	if len(w.buckets) == 0 {
		w.bucketAt = time.Time{}
		return ""
	}
	counts := make(map[string]interface{}, len(w.buckets)+2)
	for k, v := range w.buckets {
		counts[k] = v
	}
	counts["_window"] = w.duration.String()
	counts["_field"] = w.field
	out, _ := json.Marshal(counts)
	w.buckets = make(map[string]int)
	w.bucketAt = time.Time{}
	return string(out)
}

// ParseWindowFlag parses "field:duration" e.g. "level:10s".
func ParseWindowFlag(s string) (*Window, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("window: expected field:duration, got %q", s)
	}
	field := strings.TrimSpace(parts[0])
	durStr := strings.TrimSpace(parts[1])
	// Allow plain seconds as integer shorthand.
	if _, err := strconv.Atoi(durStr); err == nil {
		durStr += "s"
	}
	d, err := time.ParseDuration(durStr)
	if err != nil {
		return nil, fmt.Errorf("window: invalid duration %q: %w", durStr, err)
	}
	return NewWindow(field, d)
}
