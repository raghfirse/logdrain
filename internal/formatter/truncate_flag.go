package formatter

import "fmt"

// TruncateFlag implements flag.Value for the --truncate flag.
type TruncateFlag struct {
	MaxLen int
	Suffix string
	set    bool
}

func (f *TruncateFlag) Set(s string) error {
	maxLen, suffix, err := ParseTruncateFlag(s)
	if err != nil {
		return err
	}
	f.MaxLen = maxLen
	f.Suffix = suffix
	f.set = true
	return nil
}

func (f *TruncateFlag) String() string {
	if !f.set || f.MaxLen == 0 {
		return "0"
	}
	return fmt.Sprintf("%d:%s", f.MaxLen, f.Suffix)
}

func (f *TruncateFlag) Type() string { return "int[:suffix]" }

// Truncator returns a configured Truncator, or nil if disabled.
func (f *TruncateFlag) Truncator() *Truncator {
	if f.MaxLen <= 0 {
		return nil
	}
	return NewTruncator(f.MaxLen, f.Suffix)
}
