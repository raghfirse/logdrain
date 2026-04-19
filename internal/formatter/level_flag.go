package formatter

import "fmt"

// LevelFlag is a pflag-compatible value for --level.
type LevelFlag struct {
	Value string
}

func (f *LevelFlag) String() string {
	if f.Value == "" {
		return "info"
	}
	return f.Value
}

func (f *LevelFlag) Set(s string) error {
	norm := NormalizeLevel(s)
	if _, ok := LevelOrder[norm]; !ok {
		return fmt.Errorf("unknown log level: %q (valid: trace, debug, info, warn, error, fatal)", s)
	}
	f.Value = norm
	return nil
}

func (f *LevelFlag) Type() string { return "level" }
