package formatter

import (
	"fmt"
	"strings"
)

// existsFlagValue implements flag.Value for the --exists flag.
type existsFlagValue struct {
	config ExistsConfig
	raw    string
}

// NewExistsFlagValue returns a new existsFlagValue.
func NewExistsFlagValue() *existsFlagValue {
	return &existsFlagValue{}
}

func (f *existsFlagValue) Set(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("--exists: value must not be empty")
	}
	cfg, err := ParseExistsFlag(s)
	if err != nil {
		return fmt.Errorf("--exists: %w", err)
	}
	f.config = cfg
	f.raw = s
	return nil
}

func (f *existsFlagValue) String() string {
	if f.raw == "" {
		return ""
	}
	return f.raw
}

func (f *existsFlagValue) Type() string {
	return "exists"
}

// Config returns the parsed ExistsConfig.
func (f *existsFlagValue) Config() ExistsConfig {
	return f.config
}

// ExistsFlagUsage returns the usage string for the --exists flag.
const ExistsFlagUsage = `filter lines by field presence.
Format: comma-separated field names; prefix with '!' to require absence.
Examples:
  --exists user            keep lines that have a 'user' field
  --exists '!debug'        keep lines that do NOT have a 'debug' field
  --exists 'user,!debug'   keep lines with 'user' but without 'debug'
`
