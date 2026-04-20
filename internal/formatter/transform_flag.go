package formatter

import (
	"fmt"
	"strings"
)

// TransformFlag implements flag.Value for repeated --transform flags.
// Each value is a "from:to" rename rule.
type TransformFlag struct {
	rules []string
}

// Set appends a new transform rule after validating its format.
func (f *TransformFlag) Set(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid transform rule %q: expected \"from:to\"", s)
	}
	f.rules = append(f.rules, s)
	return nil
}

// String returns a comma-separated list of rules, or "none" if empty.
func (f *TransformFlag) String() string {
	if len(f.rules) == 0 {
		return "none"
	}
	return strings.Join(f.rules, ",")
}

// Type returns the flag type name for help output.
func (f *TransformFlag) Type() string { return "from:to" }

// Rules returns the collected rule strings.
func (f *TransformFlag) Rules() []string { return f.rules }

// ParseTransformFlag builds a *Transform from a TransformFlag, returning nil
// when no rules have been set.
func ParseTransformFlag(f *TransformFlag) (*Transform, error) {
	if f == nil || len(f.rules) == 0 {
		return &Transform{}, nil
	}
	return NewTransform(f.rules)
}
