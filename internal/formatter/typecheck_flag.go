package formatter

import (
	"fmt"
	"strings"
)

// TypeCheckFlagValue implements flag.Value for --typecheck flags.
// Multiple --typecheck flags may be provided; each is parsed independently.
type TypeCheckFlagValue struct {
	checkers []*TypeChecker
	raw      []string
}

// NewTypeCheckFlagValue returns an empty TypeCheckFlagValue.
func NewTypeCheckFlagValue() *TypeCheckFlagValue {
	return &TypeCheckFlagValue{}
}

// Set parses one flag value and appends the resulting TypeChecker.
func (f *TypeCheckFlagValue) Set(s string) error {
	c, err := ParseTypeCheckFlag(s)
	if err != nil {
		return err
	}
	f.checkers = append(f.checkers, c)
	f.raw = append(f.raw, s)
	return nil
}

// String returns the joined raw flag values, or "none" if unset.
func (f *TypeCheckFlagValue) String() string {
	if len(f.raw) == 0 {
		return "none"
	}
	return strings.Join(f.raw, ", ")
}

// Type satisfies the pflag Value interface.
func (f *TypeCheckFlagValue) Type() string { return "typecheck" }

// Checkers returns the parsed TypeChecker slice.
func (f *TypeCheckFlagValue) Checkers() []*TypeChecker { return f.checkers }

// TypeCheckFlagUsage returns the usage string for --typecheck.
func TypeCheckFlagUsage() string {
	return fmt.Sprintf(
		"filter lines by field JSON type (key:type or !key:type to invert); " +
			"type one of: string, number, bool, array, object, null (repeatable)",
	)
}
