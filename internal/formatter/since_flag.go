package formatter

import (
	"fmt"
	"time"
)

// SinceFlag is a flag.Value implementation for --since.
type SinceFlag struct {
	raw    string
	Cutoff time.Time
}

// Set parses and stores a since duration string.
func (f *SinceFlag) Set(s string) error {
	t, err := ParseSinceFlag(s)
	if err != nil {
		return err
	}
	f.raw = s
	f.Cutoff = t
	return nil
}

// String returns the raw flag value.
func (f *SinceFlag) String() string {
	if f.raw == "" {
		return ""
	}
	return f.raw
}

// Type returns the flag type name for help text.
func (f *SinceFlag) Type() string {
	return "duration"
}

// IsSet reports whether the flag has been explicitly set.
func (f *SinceFlag) IsSet() bool {
	return f.raw != ""
}

// SinceFlagUsage is the usage string for the --since flag.
const SinceFlagUsage = `only show log lines with a timestamp within this duration ago
(e.g. "5m", "1h", "30s"). Requires a parseable timestamp field.`

func (f *SinceFlag) GoString() string {
	if f.raw == "" {
		return "SinceFlag{}"
	}
	return fmt.Sprintf("SinceFlag{%q, %v}", f.raw, f.Cutoff)
}
