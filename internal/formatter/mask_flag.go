package formatter

import (
	"fmt"
	"strings"
)

// maskFlagValue implements flag.Value for repeated --mask flags.
type maskFlagValue struct {
	specs *[]string
}

func (f *maskFlagValue) Set(val string) error {
	val = strings.TrimSpace(val)
	if val == "" {
		return fmt.Errorf("mask: spec must not be empty")
	}
	parts := strings.SplitN(val, ":", 3)
	if len(parts) != 3 {
		return fmt.Errorf("mask: invalid spec %q: expected field:pattern:replacement", val)
	}
	*f.specs = append(*f.specs, val)
	return nil
}

func (f *maskFlagValue) String() string {
	if f.specs == nil || len(*f.specs) == 0 {
		return ""
	}
	return strings.Join(*f.specs, ", ")
}

func (f *maskFlagValue) Type() string { return "field:pattern:replacement" }

// NewMaskFlagValue returns a flag.Value backed by the given specs slice.
func NewMaskFlagValue(specs *[]string) *maskFlagValue {
	return &maskFlagValue{specs: specs}
}

// MaskFlagUsage is the help text for the --mask flag.
const MaskFlagUsage = `mask a field value using a regex replacement (repeatable).
Format: field:pattern:replacement
Example: --mask 'card:[0-9]{12}(\d{4}):****$1'`
