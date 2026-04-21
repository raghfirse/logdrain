package formatter

import (
	"fmt"
)

// FlattenFlag implements flag.Value for the --flatten flag.
type FlattenFlag struct {
	Value     *Flattener
	set       bool
	separator string
}

func (f *FlattenFlag) Set(s string) error {
	v, err := ParseFlattenFlag(s)
	if err != nil {
		return err
	}
	f.Value = v
	f.separator = s
	if f.separator == "" {
		f.separator = "."
	}
	f.set = true
	return nil
}

func (f *FlattenFlag) String() string {
	if !f.set {
		return ""
	}
	return fmt.Sprintf("flatten(sep=%q)", f.separator)
}

func (f *FlattenFlag) Type() string {
	return "separator"
}

// IsSet reports whether the flag was explicitly set.
func (f *FlattenFlag) IsSet() bool {
	return f.set
}
