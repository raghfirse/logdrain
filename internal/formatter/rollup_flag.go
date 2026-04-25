package formatter

import (
	"fmt"
	"strings"
)

// RollupFlagValue is a flag.Value implementation for the --rollup flag.
type RollupFlagValue struct {
	field string
	set   bool
}

// Set parses and validates the rollup field name.
func (f *RollupFlagValue) Set(s string) error {
	v, err := ParseRollupFlag(s)
	if err != nil {
		return err
	}
	f.field = v
	f.set = true
	return nil
}

// String returns the current field name or empty string if unset.
func (f *RollupFlagValue) String() string {
	if !f.set {
		return ""
	}
	return f.field
}

// Type returns the flag type name for help output.
func (f *RollupFlagValue) Type() string {
	return "field"
}

// Field returns the parsed field name.
func (f *RollupFlagValue) Field() string {
	return f.field
}

// IsSet reports whether the flag was explicitly set.
func (f *RollupFlagValue) IsSet() bool {
	return f.set
}

// RollupFlagUsage is the usage string for the --rollup flag.
const RollupFlagUsage = strings.TrimSpace(`
Group log lines by a JSON field and emit a count summary on flush.
Example: --rollup service
`)

func init() {
	// Ensure the const compiles; referenced by cmd layer.
	_ = fmt.Sprintf("%s", RollupFlagUsage)
}
