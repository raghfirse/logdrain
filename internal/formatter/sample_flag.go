package formatter

import (
	"fmt"
	"strconv"
)

// ParseSampleFlag parses a --sample flag value into a uint64 rate.
// A value of "0" or "1" disables sampling.
func ParseSampleFlag(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid sample rate %q: must be a non-negative integer", s)
	}
	return v, nil
}

// SampleFlag implements flag.Value for --sample.
type SampleFlag struct {
	Value uint64
}

func (f *SampleFlag) Set(s string) error {
	v, err := ParseSampleFlag(s)
	if err != nil {
		return err
	}
	f.Value = v
	return nil
}

func (f *SampleFlag) String() string {
	if f.Value == 0 {
		return "0"
	}
	return strconv.FormatUint(f.Value, 10)
}

func (f *SampleFlag) Type() string { return "uint" }
