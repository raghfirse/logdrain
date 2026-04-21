package formatter

import (
	"fmt"
	"strings"
)

// JoinFlagValue implements flag.Value for --join flags.
type JoinFlagValue struct {
	rules []JoinRule
}

// Set parses and appends a join rule.
func (f *JoinFlagValue) Set(s string) error {
	rule, err := ParseJoinFlag(s)
	if err != nil {
		return err
	}
	f.rules = append(f.rules, rule)
	return nil
}

// String returns a human-readable representation of all join rules.
func (f *JoinFlagValue) String() string {
	if len(f.rules) == 0 {
		return ""
	}
	parts := make([]string, len(f.rules))
	for i, r := range f.rules {
		parts[i] = fmt.Sprintf("%s:%s->%s", strings.Join(r.Keys, "+"), r.Sep, r.OutputKey)
	}
	return strings.Join(parts, ", ")
}

// Type returns the flag type name.
func (f *JoinFlagValue) Type() string { return "join" }

// Rules returns the parsed join rules.
func (f *JoinFlagValue) Rules() []JoinRule { return f.rules }
