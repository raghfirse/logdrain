package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	gotemplate "text/template"
)

// Template renders a log line using a Go text/template string.
// Fields from the parsed JSON object are available as .Field (dot-access).
// If the line is not valid JSON, it is passed through unchanged.
type Template struct {
	tmpl *gotemplate.Template
}

// NewTemplate compiles the given template string and returns a Template.
// Returns an error if the template cannot be parsed.
func NewTemplate(tmplStr string) (*Template, error) {
	if strings.TrimSpace(tmplStr) == "" {
		return nil, fmt.Errorf("template string must not be empty")
	}
	t, err := gotemplate.New("logline").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}
	return &Template{tmpl: t}, nil
}

// Render applies the template to the given log line.
// The line is decoded as a JSON object; its keys are available by name.
// If decoding fails, the original line is returned unchanged.
func (t *Template) Render(line string) string {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return line
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, fields); err != nil {
		return line
	}
	return buf.String()
}

// ParseTemplateFlag parses a --template flag value, returning a *Template
// or an error if the value is invalid.
func ParseTemplateFlag(value string) (*Template, error) {
	return NewTemplate(value)
}
