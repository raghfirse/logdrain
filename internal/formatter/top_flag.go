package formatter

import (
	"fmt"
	"strings"
)

// TopFlagValue implements flag.Value for --top.
type TopFlagValue struct {
	Field string
	N     int
	set   bool
}

// Set parses "field:n" and stores the result.
func (f *TopFlagValue) Set(s string) error {
	field, n, err := ParseTopFlag(s)
	if err != nil {
		return err
	}
	if _, err := NewTop(field, n); err != nil {
		return err
	}
	f.Field = field
	f.N = n
	f.set = true
	return nil
}

// String returns the current value as a string.
func (f *TopFlagValue) String() string {
	if !f.set {
		return ""
	}
	return fmt.Sprintf("%s:%d", f.Field, f.N)
}

// Type returns the flag type name.
func (f *TopFlagValue) Type() string {
	return "field:n"
}

// IsSet reports whether the flag was explicitly set.
func (f *TopFlagValue) IsSet() bool {
	return f.set
}

// TopFlagUsage returns the usage string for --top.
func TopFlagUsage() string {
	return strings.TrimSpace(`
Report the top N most frequent values for a JSON field.
Format: field:n (e.g. --top level:5)
Output is emitted at program exit via Flush.
`)
}
