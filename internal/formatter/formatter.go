package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Format controls how log lines are rendered.
type Format int

const (
	FormatPretty Format = iota
	FormatRaw
)

// Formatter writes formatted log entries to an output writer.
type Formatter struct {
	format Format
	out    io.Writer
}

// New creates a Formatter writing to out with the given format.
func New(out io.Writer, format Format) *Formatter {
	return &Formatter{format: format, out: out}
}

// Write formats and writes a single log line, prefixed with the source name.
func (f *Formatter) Write(source, line string) error {
	if f.format == FormatRaw {
		_, err := fmt.Fprintf(f.out, "[%s] %s\n", source, line)
		return err
	}
	return f.writePretty(source, line)
}

func (f *Formatter) writePretty(source, line string) error {
	var entry map[string]any
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		// Not valid JSON — fall back to raw.
		_, err = fmt.Fprintf(f.out, "[%s] %s\n", source, line)
		return err
	}

	level := strings.ToUpper(stringField(entry, "level", "INFO"))
	msg := stringField(entry, "msg", stringField(entry, "message", line))
	ts := stringField(entry, "time", stringField(entry, "timestamp", ""))

	prefix := ""
	if ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			prefix = t.Format("15:04:05") + " "
		} else {
			prefix = ts + " "
		}
	}

	_, err := fmt.Fprintf(f.out, "%s[%s] %-5s %s\n", prefix, source, level, msg)
	return err
}

func stringField(m map[string]any, key, fallback string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return fallback
}
