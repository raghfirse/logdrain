package formatter

import (
	"fmt"
	"strings"
)

// ClampFlagValue implements flag.Value for --clamp flags.
type ClampFlagValue struct {
	clamps []*Clamp
	raw    []string
}

// NewClampFlagValue returns an initialised ClampFlagValue.
func NewClampFlagValue() *ClampFlagValue {
	return &ClampFlagValue{}
}

// Set parses and appends a clamp spec (field,min,max).
func (f *ClampFlagValue) Set(s string) error {
	c, err := ParseClampFlag(s)
	if err != nil {
		return err
	}
	f.clamps = append(f.clamps, c)
	f.raw = append(f.raw, s)
	return nil
}

// String returns the current flag value as a human-readable string.
func (f *ClampFlagValue) String() string {
	if len(f.raw) == 0 {
		return ""
	}
	return strings.Join(f.raw, "; ")
}

// Type returns the flag type name for help output.
func (f *ClampFlagValue) Type() string {
	return "field,min,max"
}

// Clamps returns the parsed Clamp transformers.
func (f *ClampFlagValue) Clamps() []*Clamp {
	return f.clamps
}

// ClampFlagUsage is the usage string shown in --help output.
const ClampFlagUsage = `Clamp a numeric field to [min,max]. Format: field,min,max.
May be specified multiple times.
Example: --clamp latency,0,5000`

func clampFlagHelp() string {
	return fmt.Sprintf("  --clamp %s\n    \t%s", "field,min,max", ClampFlagUsage)
}
