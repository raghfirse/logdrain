package formatter

import (
	"fmt"
	"strings"
)

// DiffFlagValue implements flag.Value for the --diff flag.
type DiffFlagValue struct {
	fields []string
	raw    string
}

// Set parses a comma-separated list of field names.
func (f *DiffFlagValue) Set(s string) error {
	fields, err := ParseDiffFlag(s)
	if err != nil {
		return err
	}
	f.fields = fields
	f.raw = s
	return nil
}

// String returns the current flag value as a string.
func (f *DiffFlagValue) String() string {
	if f.raw == "" {
		return ""
	}
	return f.raw
}

// Type returns the flag type name for help output.
func (f *DiffFlagValue) Type() string {
	return "fields"
}

// Fields returns the parsed field names.
func (f *DiffFlagValue) Fields() []string {
	return f.fields
}

// DiffFlagUsage returns the usage string for the --diff flag.
func DiffFlagUsage() string {
	lines := []string{
		"emit only changed fields between consecutive JSON lines",
		"  (empty = diff all fields)",
		"  examples:",
		"    --diff              diff every field",
		"    --diff=status,code  diff only 'status' and 'code'",
	}
	return strings.Join(lines, "\n")
}

// NewDiffFlagValue constructs a DiffFlagValue with optional initial fields.
func NewDiffFlagValue(fields ...string) *DiffFlagValue {
	if len(fields) == 0 {
		return &DiffFlagValue{}
	}
	return &DiffFlagValue{
		fields: fields,
		raw:    fmt.Sprintf("%s", strings.Join(fields, ",")),
	}
}
