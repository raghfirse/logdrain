package formatter

import (
	"fmt"
	"strings"
)

// pickFlagValue implements flag.Value for the --pick flag.
type pickFlagValue struct {
	keys []string
}

func (f *pickFlagValue) Set(s string) error {
	keys, err := ParsePickFlag(s)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("--pick requires at least one field name")
	}
	f.keys = append(f.keys, keys...)
	return nil
}

func (f *pickFlagValue) String() string {
	if len(f.keys) == 0 {
		return ""
	}
	return strings.Join(f.keys, ",")
}

func (f *pickFlagValue) Type() string {
	return "fields"
}

func (f *pickFlagValue) Keys() []string {
	return f.keys
}

// NewPickFlagValue returns a new pickFlagValue for use with flag.Var.
func NewPickFlagValue() *pickFlagValue {
	return &pickFlagValue{}
}

// PickFlagUsage is the usage string for the --pick flag.
const PickFlagUsage = `retain only the given comma-separated fields in each JSON log line (e.g. --pick msg,level,ts)`
