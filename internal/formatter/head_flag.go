package formatter

import (
	"fmt"
)

// HeadFlagValue implements flag.Value for the --head flag.
type HeadFlagValue struct {
	head *Head
	raw  string
}

func (f *HeadFlagValue) Set(s string) error {
	h, err := ParseHeadFlag(s)
	if err != nil {
		return err
	}
	f.head = h
	f.raw = s
	return nil
}

func (f *HeadFlagValue) String() string {
	if f.raw == "" {
		return ""
	}
	return f.raw
}

func (f *HeadFlagValue) Type() string {
	return "head"
}

// Head returns the parsed Head, or nil if not set.
func (f *HeadFlagValue) Head() *Head {
	return f.head
}

// ParseHeadFlagValue parses and validates the head flag string, returning a
// descriptive error suitable for CLI output.
func ParseHeadFlagValue(s string) (*HeadFlagValue, error) {
	if s == "" {
		return &HeadFlagValue{}, nil
	}
	h, err := ParseHeadFlag(s)
	if err != nil {
		return nil, fmt.Errorf("--head: %w", err)
	}
	return &HeadFlagValue{head: h, raw: s}, nil
}
