package formatter

import (
	"fmt"
)

// WindowFlagValue implements flag.Value for --window.
type WindowFlagValue struct {
	raw string
	W   *Window
}

func (f *WindowFlagValue) Set(s string) error {
	w, err := ParseWindowFlag(s)
	if err != nil {
		return err
	}
	f.raw = s
	f.W = w
	return nil
}

func (f *WindowFlagValue) String() string {
	if f.raw == "" {
		return ""
	}
	return f.raw
}

func (f *WindowFlagValue) Type() string {
	return "field:duration"
}

// NewWindowFlagValue returns an empty WindowFlagValue suitable for flag registration.
func NewWindowFlagValue() *WindowFlagValue {
	return &WindowFlagValue{}
}

// Usage returns a help string for the --window flag.
func WindowFlagUsage() string {
	return fmt.Sprintf(
		"group log lines by field into time buckets and emit a count summary\n" +
			"  format: field:duration  (e.g. level:10s, status:1m)\n" +
			"  when the window expires a JSON summary line is emitted",
	)
}
