package formatter

import "fmt"

// TimestampMode controls how timestamps are rendered.
type TimestampMode int

const (
	TimestampOff   TimestampMode = iota // do not print timestamp
	TimestampShort                      // HH:MM:SS
	TimestampFull                       // YYYY-MM-DD HH:MM:SS.mmm
)

// ParseTimestampMode converts a flag string value to a TimestampMode.
func ParseTimestampMode(s string) (TimestampMode, error) {
	switch s {
	case "off", "none", "":
		return TimestampOff, nil
	case "short":
		return TimestampShort, nil
	case "full", "long":
		return TimestampFull, nil
	}
	return TimestampOff, fmt.Errorf("unknown timestamp mode %q: want off|short|full", s)
}

// String implements fmt.Stringer.
func (m TimestampMode) String() string {
	switch m {
	case TimestampShort:
		return "short"
	case TimestampFull:
		return "full"
	default:
		return "off"
	}
}
