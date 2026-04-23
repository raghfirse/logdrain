package formatter

import "fmt"

// SkipFlagValue implements flag.Value for the --skip flag.
type SkipFlagValue struct {
	raw    string
	Skipper *Skipper
}

// Set parses the flag value and stores the resulting Skipper.
func (f *SkipFlagValue) Set(s string) error {
	sk, err := ParseSkipFlag(s)
	if err != nil {
		return err
	}
	f.raw = s
	f.Skipper = sk
	return nil
}

// String returns the raw flag value, or "0" if unset.
func (f *SkipFlagValue) String() string {
	if f.raw == "" {
		return "0"
	}
	return f.raw
}

// Type returns the flag type name for help text.
func (f *SkipFlagValue) Type() string { return "skip" }

// SkipFlagUsage returns the usage string for the --skip flag.
func SkipFlagUsage() string {
	return fmt.Sprintf(
		"skip the first N lines before emitting output (format: N or N:field)\n" +
			"  --skip 10          skip first 10 lines globally\n" +
			"  --skip 5:service   skip first 5 lines per distinct 'service' value",
	)
}
