package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format controls output style.
type Format int

const (
	Raw    Format = iota // emit line as-is
	Pretty               // human-friendly colored output
)

// Formatter writes log lines to an io.Writer.
type Formatter struct {
	out    io.Writer
	fmt    Format
	color  bool
}

// New creates a Formatter writing to out.
func New(out io.Writer, f Format, color bool) *Formatter {
	return &Formatter{out: out, fmt: f, color: color}
}

// Write formats and writes a single log line from the named source.
func (f *Formatter) Write(source, line string) error {
	if f.fmt == Raw {
		_, err := fmt.Fprintf(f.out, "%s\n", line)
		return err
	}
	return f.writePretty(source, line)
}

func (f *Formatter) writePretty(source, line string) error {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		// Not valid JSON — emit raw with source tag.
		src := f.formatSource(source)
		_, err2 := fmt.Fprintf(f.out, "%s %s\n", src, line)
		return err2
	}

	level := stringField(fields, "level", "lvl")
	msg := stringField(fields, "message", "msg")
	ts := stringField(fields, "time", "ts", "timestamp")

	var sb strings.Builder
	sb.WriteString(f.formatSource(source))
	if ts != "" {
		sb.WriteString(" " + f.maybeColor(colorGray, ts))
	}
	if level != "" {
		sb.WriteString(" " + f.maybeColor(LevelColor(level), strings.ToUpper(level)))
	}
	if msg != "" {
		sb.WriteString(" " + msg)
	}

	// Append remaining fields sorted.
	skip := map[string]bool{"level": true, "lvl": true, "message": true, "msg": true, "time": true, "ts": true, "timestamp": true}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		if !skip[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(" %s=%v", f.maybeColor(colorGray, k), fields[k]))
	}

	_, err := fmt.Fprintln(f.out, sb.String())
	return err
}

func (f *Formatter) formatSource(source string) string {
	tag := fmt.Sprintf("[%s]", source)
	if f.color {
		return Colorize(SourceColor(source), tag)
	}
	return tag
}

func (f *Formatter) maybeColor(color, text string) string {
	if f.color {
		return Colorize(color, text)
	}
	return text
}

func stringField(fields map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := fields[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}
