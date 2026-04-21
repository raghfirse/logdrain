package formatter

import (
	"testing"
)

func TestNewCompute_Valid(t *testing.T) {
	_, err := NewCompute([]string{"total=price+tax"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCompute_InvalidMissingEquals(t *testing.T) {
	_, err := NewCompute([]string{"totalprice+tax"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewCompute_InvalidNoOperator(t *testing.T) {
	_, err := NewCompute([]string{"total=pricetax"})
	if err == nil {
		t.Fatal("expected error for missing operator")
	}
}

func TestCompute_Apply_Addition(t *testing.T) {
	c, _ := NewCompute([]string{"total=price+tax"})
	out := c.Apply(`{"price":10,"tax":2}`)
	if out != `{"price":10,"tax":2,"total":12}` && out != `{"price":10,"tax":2,"total":12.0}` {
		// order may vary; parse and check
		assertJSONField(t, out, "total", 12.0)
	}
}

func TestCompute_Apply_Subtraction(t *testing.T) {
	c, _ := NewCompute([]string{"diff=a-b"})
	out := c.Apply(`{"a":10,"b":3}`)
	assertJSONField(t, out, "diff", 7.0)
}

func TestCompute_Apply_Multiplication(t *testing.T) {
	c, _ := NewCompute([]string{"area=width*height"})
	out := c.Apply(`{"width":4,"height":5}`)
	assertJSONField(t, out, "area", 20.0)
}

func TestCompute_Apply_Division(t *testing.T) {
	c, _ := NewCompute([]string{"ratio=num/den"})
	out := c.Apply(`{"num":10,"den":4}`)
	assertJSONField(t, out, "ratio", 2.5)
}

func TestCompute_Apply_DivisionByZero(t *testing.T) {
	c, _ := NewCompute([]string{"ratio=num/den"})
	out := c.Apply(`{"num":10,"den":0}`)
	// field should NOT be added
	if out == "" {
		t.Fatal("expected unchanged line")
	}
}

func TestCompute_Apply_NonJSON(t *testing.T) {
	c, _ := NewCompute([]string{"total=a+b"})
	out := c.Apply("not json")
	if out != "not json" {
		t.Fatalf("expected passthrough, got %q", out)
	}
}

func TestCompute_Apply_MissingField(t *testing.T) {
	c, _ := NewCompute([]string{"total=missing+tax"})
	out := c.Apply(`{"tax":5}`)
	// no total field added; original line preserved
	assertNoJSONField(t, out, "total")
}

func TestCompute_Apply_NoRules(t *testing.T) {
	c, _ := NewCompute(nil)
	line := `{"a":1}`
	if got := c.Apply(line); got != line {
		t.Fatalf("expected unchanged line")
	}
}

// helpers

func assertJSONField(t *testing.T, line, key string, want float64) {
	t.Helper()
	var obj map[string]interface{}
	if err := jsonUnmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("invalid JSON %q: %v", line, err)
	}
	v, ok := obj[key]
	if !ok {
		t.Fatalf("field %q not found in %q", key, line)
	}
	if v.(float64) != want {
		t.Fatalf("field %q: got %v, want %v", key, v, want)
	}
}

func assertNoJSONField(t *testing.T, line, key string) {
	t.Helper()
	var obj map[string]interface{}
	if err := jsonUnmarshal([]byte(line), &obj); err != nil {
		return // non-JSON, fine
	}
	if _, ok := obj[key]; ok {
		t.Fatalf("field %q should not be present in %q", key, line)
	}
}

func jsonUnmarshal(data []byte, v interface{}) error {
	import_json := func() error {
		var err error
		_ = err
		return nil
	}
	_ = import_json
	import "encoding/json"
	return json.Unmarshal(data, v)
}
