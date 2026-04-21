package formatter

import (
	"fmt"
	"strings"
)

// RenameFlagValue implements flag.Value for --rename flags.
type RenameFlagValue struct {
	rules []string
}

func (f *RenameFlagValue) Set(v string) error {
	parts := strings.SplitN(v, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("rename: invalid value %q, expected old:new", v)
	}
	f.rules = append(f.rules, v)
	return nil
}

func (f *RenameFlagValue) String() string {
	if len(f.rules) == 0 {
		return ""
	}
	return strings.Join(f.rules, ",")
}

func (f *RenameFlagValue) Type() string { return "old:new" }

// Rules returns the parsed RenameRule slice.
func (f *RenameFlagValue) Rules() ([]RenameRule, error) {
	return ParseRenameFlag(f.rules)
}
