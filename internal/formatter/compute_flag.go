package formatter

import (
	"fmt"
	"strings"
)

// computeFlagValue implements flag.Value for repeated --compute flags.
type computeFlagValue struct {
	exprs []string
}

func (f *computeFlagValue) Set(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("compute expression must not be empty")
	}
	// validate eagerly
	if _, err := parseComputeExpr(s); err != nil {
		return err
	}
	f.exprs = append(f.exprs, s)
	return nil
}

func (f *computeFlagValue) String() string {
	if len(f.exprs) == 0 {
		return ""
	}
	return strings.Join(f.exprs, ", ")
}

func (f *computeFlagValue) Type() string { return "expr" }

// ParseComputeFlag builds a *Compute from a computeFlagValue.
func ParseComputeFlag(f *computeFlagValue) (*Compute, error) {
	if f == nil || len(f.exprs) == 0 {
		return NewCompute(nil)
	}
	return NewCompute(f.exprs)
}
