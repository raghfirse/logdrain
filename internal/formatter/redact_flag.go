package formatter

import "strings"

// RedactFlag implements flag.Value for --redact.
type RedactFlag struct {
	values []string
}

func (f *RedactFlag) Set(s string) error {
	keys := ParseRedactFlag(s)
	f.values = append(f.values, keys...)
	return nil
}

func (f *RedactFlag) String() string {
	if len(f.values) == 0 {
		return ""
	}
	return strings.Join(f.values, ",")
}

func (f *RedactFlag) Type() string {
	return "fields"
}

// Keys returns the accumulated list of field names to redact.
func (f *RedactFlag) Keys() []string {
	return f.values
}
