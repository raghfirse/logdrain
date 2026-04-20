package formatter

import "fmt"

// TemplateFlag is a flag.Value implementation for --template.
type TemplateFlag struct {
	value string
	tmpl  *Template
}

// Set parses and compiles the template string.
func (f *TemplateFlag) Set(s string) error {
	t, err := ParseTemplateFlag(s)
	if err != nil {
		return err
	}
	f.value = s
	f.tmpl = t
	return nil
}

// String returns the current template string, or an empty string if unset.
func (f *TemplateFlag) String() string {
	return f.value
}

// Type returns the flag type name for CLI help output.
func (f *TemplateFlag) Type() string {
	return "template"
}

// Template returns the compiled *Template, or nil if not set.
func (f *TemplateFlag) Template() *Template {
	return f.tmpl
}

// IsSet reports whether a template string has been provided.
func (f *TemplateFlag) IsSet() bool {
	return f.tmpl != nil
}

// ParseTemplateFlagValue is a convenience wrapper used by main to validate
// and return a formatted error message suitable for CLI output.
func ParseTemplateFlagValue(s string) (*Template, error) {
	t, err := ParseTemplateFlag(s)
	if err != nil {
		return nil, fmt.Errorf("--template: %w", err)
	}
	return t, nil
}
