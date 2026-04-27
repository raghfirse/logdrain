package formatter

import (
	"fmt"
	"strings"
)

// OmitFlagValue implements flag.Value for --omit.
type OmitFlagValue struct {
	keys []string
}

func (f *OmitFlagValue) Set(s string) error {
	keys, err := ParseOmitFlag(s)
	if err != nil {
		return err
	}
	f.keys = append(f.keys, keys...)
	return nil
}

func (f *OmitFlagValue) String() string {
	if len(f.keys) == 0 {
		return ""
	}
	return strings.Join(f.keys, ",")
}

func (f *OmitFlagValue) Type() string {
	return "key[,key...]"
}

// Keys returns the accumulated list of keys to omit.
func (f *OmitFlagValue) Keys() []string {
	return f.keys
}

// NewOmitFlagValue returns an initialised OmitFlagValue.
func NewOmitFlagValue() *OmitFlagValue {
	return &OmitFlagValue{}
}

// OmitFlagUsage is the help text shown for --omit.
const OmitFlagUsage = `Remove one or more keys from every JSON log line.
Accepts a comma-separated list of field names (case-insensitive).
Example: --omit secret,password
         --omit Authorization
`

func init() {
	_ = fmt.Sprintf // keep import
}
