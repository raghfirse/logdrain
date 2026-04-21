package formatter

import (
	"fmt"
	"strconv"
)

// LimitFlagValue implements flag.Value for --limit.
type LimitFlagValue struct {
	N int
}

// Set parses the flag value.
func (f *LimitFlagValue) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("limit: invalid integer %q", s)
	}
	if n < 0 {
		return fmt.Errorf("limit: value must be >= 0, got %d", n)
	}
	f.N = n
	return nil
}

// String returns the current value as a string.
func (f *LimitFlagValue) String() string {
	if f.N == 0 {
		return "0"
	}
	return strconv.Itoa(f.N)
}

// Type returns the flag type name used in help output.
func (f *LimitFlagValue) Type() string {
	return "int"
}

// LimitFlagUsage is the usage string for --limit.
const LimitFlagUsage = "stop after emitting N lines (0 = unlimited)"
