package formatter

import (
	"fmt"
	"strings"
)

// ReplaceFlagValue implements flag.Value for --replace flags.
type ReplaceFlagValue struct {
	rules []ReplaceRule
}

// NewReplaceFlagValue returns an empty ReplaceFlagValue.
func NewReplaceFlagValue() *ReplaceFlagValue {
	return &ReplaceFlagValue{}
}

// Set parses and appends a replace rule in "key/pattern/replacement" form.
func (f *ReplaceFlagValue) Set(s string) error {
	rule, err := ParseReplaceFlag(s)
	if err != nil {
		return err
	}
	f.rules = append(f.rules, rule)
	return nil
}

// String returns a human-readable representation of the accumulated rules.
func (f *ReplaceFlagValue) String() string {
	if len(f.rules) == 0 {
		return ""
	}
	parts := make([]string, len(f.rules))
	for i, r := range f.rules {
		parts[i] = fmt.Sprintf("%s/%s/%s", r.Key, r.Pattern.String(), r.With)
	}
	return strings.Join(parts, ", ")
}

// Type returns the flag type name.
func (f *ReplaceFlagValue) Type() string { return "key/pattern/replacement" }

// Rules returns the parsed replace rules.
func (f *ReplaceFlagValue) Rules() []ReplaceRule { return f.rules }

// ReplaceFlagUsage is the help text for the --replace flag.
const ReplaceFlagUsage = `Replace field values using a regex (repeatable).
Format: key/pattern/replacement
Example: --replace message/foo/bar  replaces "foo" with "bar" in the message field.`
