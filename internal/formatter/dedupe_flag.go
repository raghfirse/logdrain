package formatter

import (
	"fmt"
	"time"
)

// DedupeFlag is a flag value for the --dedupe duration flag.
type DedupeFlag struct {
	Duration time.Duration
	Set      bool
}

func (f *DedupeFlag) String() string {
	if !f.Set {
		return "0s"
	}
	return f.Duration.String()
}

func (f *DedupeFlag) Type() string {
	return "duration"
}

func (f *DedupeFlag) IsBoolFlag() bool { return false }

func (f *DedupeFlag) Value() *DedupeFilter {
	if !f.Set || f.Duration <= 0 {
		return nil
	}
	return NewDedupeFilter(f.Duration)
}

// ParseDedupeFlag parses a duration string into a DedupeFlag.
func ParseDedupeFlag(s string) (DedupeFlag, error) {
	if s == "" || s == "0" || s == "0s" {
		return DedupeFlag{}, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return DedupeFlag{}, fmt.Errorf("invalid dedupe duration %q: %w", s, err)
	}
	if d < 0 {
		return DedupeFlag{}, fmt.Errorf("dedupe duration must be non-negative")
	}
	return DedupeFlag{Duration: d, Set: true}, nil
}
