package filter

import (
	"encoding/json"
	"strings"
)

// Rule represents a single filter condition on a JSON log entry.
type Rule struct {
	Field string
	Value string
}

// Filter holds a set of rules and applies them to log lines.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of "field=value" expressions.
func New(exprs []string) (*Filter, error) {
	rules := make([]Rule, 0, len(exprs))
	for _, expr := range exprs {
		parts := strings.SplitN(expr, "=", 2)
		if len(parts) != 2 {
			return nil, &InvalidExprError{Expr: expr}
		}
		rules = append(rules, Rule{Field: parts[0], Value: parts[1]})
	}
	return &Filter{rules: rules}, nil
}

// Match returns true if the JSON line satisfies all filter rules.
func (f *Filter) Match(line []byte) bool {
	if len(f.rules) == 0 {
		return true
	}
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return false
	}
	for _, rule := range f.rules {
		val, ok := entry[rule.Field]
		if !ok {
			return false
		}
		if !strings.EqualFold(fmt.Sprintf("%v", val), rule.Value) {
			return false
		}
	}
	return true
}

// InvalidExprError is returned when a filter expression is malformed.
type InvalidExprError struct {
	Expr string
}

func (e *InvalidExprError) Error() string {
	return "invalid filter expression: " + e.Expr
}
