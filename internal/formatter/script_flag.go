package formatter

import (
	"fmt"
	"strings"
)

// ScriptFlag implements flag.Value for repeated --script flags.
// Each value is a "key=template" rule string.
type ScriptFlag struct {
	rules []string
}

func (f *ScriptFlag) Set(val string) error {
	val = strings.TrimSpace(val)
	if val == "" {
		return fmt.Errorf("script rule must not be empty")
	}
	// Validate eagerly so the user gets an error at startup.
	if _, err := NewScript([]string{val}); err != nil {
		return err
	}
	f.rules = append(f.rules, val)
	return nil
}

func (f *ScriptFlag) String() string {
	if len(f.rules) == 0 {
		return ""
	}
	return strings.Join(f.rules, ", ")
}

func (f *ScriptFlag) Type() string { return "key=template" }

// Rules returns the accumulated rule strings.
func (f *ScriptFlag) Rules() []string { return f.rules }

// ParseScriptFlag builds a Script from a ScriptFlag, returning nil when no
// rules have been registered.
func ParseScriptFlag(f *ScriptFlag) (*Script, error) {
	if f == nil || len(f.rules) == 0 {
		return nil, nil
	}
	return NewScript(f.rules)
}
