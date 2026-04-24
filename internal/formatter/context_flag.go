package formatter

import "strings"

// ContextFlag is a flag.Value implementation for --context.
type ContextFlag struct {
	fields []string
}

// Set parses and stores comma-separated context field names.
func (f *ContextFlag) Set(s string) error {
	f.fields = ParseContextFlag(s)
	return nil
}

// String returns the joined field names.
func (f *ContextFlag) String() string {
	return strings.Join(f.fields, ",")
}

// Type returns the flag type name.
func (f *ContextFlag) Type() string {
	return "fields"
}

// Fields returns the parsed list of field names.
func (f *ContextFlag) Fields() []string {
	return f.fields
}

// Contains reports whether the given field name is present in the context fields.
func (f *ContextFlag) Contains(field string) bool {
	for _, v := range f.fields {
		if v == field {
			return true
		}
	}
	return false
}
