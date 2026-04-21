package formatter

import (
	"fmt"
)

// mergeFlagValue implements flag.Value for --merge.
type mergeFlagValue struct {
	merger    *Merger
	raw       string
	overwrite bool
}

func (f *mergeFlagValue) Set(s string) error {
	m, err := NewMerge(s, f.overwrite)
	if err != nil {
		return err
	}
	f.merger = m
	f.raw = s
	return nil
}

func (f *mergeFlagValue) String() string {
	if f.raw == "" {
		return ""
	}
	return f.raw
}

func (f *mergeFlagValue) Type() string {
	return "json-object"
}

func (f *mergeFlagValue) Merger() *Merger {
	return f.merger
}

// NewMergeFlagValue returns a flag value for collecting merge fields.
// overwrite controls whether merged fields replace existing keys.
func NewMergeFlagValue(overwrite bool) *mergeFlagValue {
	return &mergeFlagValue{overwrite: overwrite}
}

// MergeFlagUsage returns the usage string for the --merge flag.
func MergeFlagUsage(overwrite bool) string {
	if overwrite {
		return fmt.Sprintf("JSON object whose fields are merged into every log line (overwrites existing keys)")
	}
	return "JSON object whose fields are merged into every log line (skips existing keys)"
}
