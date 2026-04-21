package formatter

import (
	"fmt"
	"strings"
)

// GrepFlag holds a compiled Grep filter built from CLI flag values.
type GrepFlag struct {
	pattern string
	fields  []string
	grep    *Grep
}

// Set parses a grep flag value of the form "pattern" or "pattern:field1,field2".
func (f *GrepFlag) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	pattern := parts[0]
	var fields []string
	if len(parts) == 2 && parts[1] != "" {
		for _, field := range strings.Split(parts[1], ",") {
			field = strings.TrimSpace(field)
			if field != "" {
				fields = append(fields, field)
			}
		}
	}
	g, err := NewGrep(pattern, fields)
	if err != nil {
		return fmt.Errorf("invalid grep pattern %q: %w", pattern, err)
	}
	f.pattern = pattern
	f.fields = fields
	f.grep = g
	return nil
}

// String returns the flag's current string representation.
func (f *GrepFlag) String() string {
	if f.pattern == "" {
		return ""
	}
	if len(f.fields) == 0 {
		return f.pattern
	}
	return f.pattern + ":" + strings.Join(f.fields, ",")
}

// Type returns the flag type name for pflag/flag compatibility.
func (f *GrepFlag) Type() string { return "grep" }

// Grep returns the compiled Grep filter, or nil if not set.
func (f *GrepFlag) Grep() *Grep { return f.grep }

// ParseGrepFlag parses a raw flag string into a GrepFlag.
func ParseGrepFlag(value string) (*GrepFlag, error) {
	if value == "" {
		return &GrepFlag{}, nil
	}
	var f GrepFlag
	if err := f.Set(value); err != nil {
		return nil, err
	}
	return &f, nil
}
