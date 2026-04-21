package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Compute adds a new field to a JSON log line whose value is derived from
// a simple arithmetic expression referencing existing numeric fields.
// Expression format: "newfield=field1+field2" or "newfield=field1-field2"
// or "newfield=field1*field2" or "newfield=field1/field2".
type Compute struct {
	rules []computeRule
}

type computeRule struct {
	output string
	left   string
	op     byte
	right  string
}

// NewCompute parses a slice of compute expressions and returns a Compute.
func NewCompute(exprs []string) (*Compute, error) {
	rules := make([]computeRule, 0, len(exprs))
	for _, expr := range exprs {
		r, err := parseComputeExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("compute: %w", err)
		}
		rules = append(rules, r)
	}
	return &Compute{rules: rules}, nil
}

func parseComputeExpr(expr string) (computeRule, error) {
	// expect: output=left<op>right
	eqIdx := -1
	for i, c := range expr {
		if c == '=' {
			eqIdx = i
			break
		}
	}
	if eqIdx <= 0 {
		return computeRule{}, fmt.Errorf("invalid expression %q: missing output= prefix", expr)
	}
	output := expr[:eqIdx]
	body := expr[eqIdx+1:]
	ops := []byte{'+', '-', '*', '/'}
	for _, op := range ops {
		for i := 1; i < len(body)-1; i++ {
			if body[i] == op {
				return computeRule{output: output, left: body[:i], op: op, right: body[i+1:]}, nil
			}
		}
	}
	return computeRule{}, fmt.Errorf("invalid expression %q: no operator (+,-,*,/) found", expr)
}

// Apply evaluates all compute rules against a JSON line and returns the
// augmented line. Non-JSON lines are returned unchanged.
func (c *Compute) Apply(line string) string {
	if len(c.rules) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range c.rules {
		lv, lok := numericField(obj, r.left)
		rv, rok := numericField(obj, r.right)
		if !lok || !rok {
			continue
		}
		var result float64
		switch r.op {
		case '+':
			result = lv + rv
		case '-':
			result = lv - rv
		case '*':
			result = lv * rv
		case '/':
			if rv == 0 {
				continue
			}
			result = lv / rv
		}
		obj[r.output] = json.RawMessage(strconv.FormatFloat(result, 'f', -1, 64))
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

func numericField(obj map[string]json.RawMessage, key string) (float64, bool) {
	raw, ok := findKeyCaseInsensitive(obj, key)
	if !ok {
		return 0, false
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return 0, false
	}
	return f, true
}
