package formatter

import (
	"testing"
)

func TestNewClamp_Valid(t *testing.T) {
	c, err := NewClamp("latency", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Clamp")
	}
}

func TestNewClamp_EmptyKey(t *testing.T) {
	_, err := NewClamp("", 0, 100)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewClamp_MinGreaterThanMax(t *testing.T) {
	_, err := NewClamp("x", 10, 5)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestClamp_Apply_NonJSON(t *testing.T) {
	c, _ := NewClamp("v", 0, 10)
	out, keep := c.Apply("not json")
	if !keep || out != "not json" {
		t.Errorf("expected passthrough, got %q %v", out, keep)
	}
}

func TestClamp_Apply_MissingField(t *testing.T) {
	c, _ := NewClamp("score", 0, 10)
	line := `{"msg":"hello"}`
	out, keep := c.Apply(line)
	if !keep || out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestClamp_Apply_ValueWithinRange(t *testing.T) {
	c, _ := NewClamp("score", 0, 100)
	line := `{"score":50}`
	out, keep := c.Apply(line)
	if !keep {
		t.Fatal("expected keep=true")
	}
	if out != `{"score":50}` {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestClamp_Apply_ValueBelowMin(t *testing.T) {
	c, _ := NewClamp("score", 0, 100)
	out, keep := c.Apply(`{"score":-5}`)
	if !keep {
		t.Fatal("expected keep=true")
	}
	if out != `{"score":0}` {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestClamp_Apply_ValueAboveMax(t *testing.T) {
	c, _ := NewClamp("score", 0, 100)
	out, keep := c.Apply(`{"score":200}`)
	if !keep {
		t.Fatal("expected keep=true")
	}
	if out != `{"score":100}` {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestClamp_Apply_NonNumericField(t *testing.T) {
	c, _ := NewClamp("score", 0, 100)
	line := `{"score":"high"}`
	out, keep := c.Apply(line)
	if !keep || out != line {
		t.Errorf("expected passthrough for non-numeric, got %q", out)
	}
}

func TestParseClampFlag_Valid(t *testing.T) {
	c, err := ParseClampFlag("latency,0,500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.key != "latency" || c.min != 0 || c.max != 500 {
		t.Errorf("unexpected clamp: %+v", c)
	}
}

func TestParseClampFlag_Invalid(t *testing.T) {
	cases := []string{
		"nocommas",
		"field,abc,100",
		"field,0,xyz",
		"field,50,10",
	}
	for _, tc := range cases {
		_, err := ParseClampFlag(tc)
		if err == nil {
			t.Errorf("expected error for %q", tc)
		}
	}
}
