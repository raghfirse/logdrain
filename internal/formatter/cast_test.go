package formatter

import (
	"testing"
)

func TestNewCast_Valid(t *testing.T) {
	_, err := NewCast([]CastRule{{Key: "count", Type: "int"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCast_EmptyKey(t *testing.T) {
	_, err := NewCast([]CastRule{{Key: "", Type: "int"}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewCast_UnknownType(t *testing.T) {
	_, err := NewCast([]CastRule{{Key: "x", Type: "bytes"}})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestCast_Apply_NonJSON(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "x", Type: "int"}})
	input := "not json at all"
	if got := c.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestCast_Apply_NoRules(t *testing.T) {
	c, _ := NewCast(nil)
	input := `{"count":"42"}`
	if got := c.Apply(input); got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestCast_Apply_StringToInt(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "count", Type: "int"}})
	got := c.Apply(`{"count":"7"}`)
	if got != `{"count":7}` {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestCast_Apply_NumberToString(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "score", Type: "string"}})
	got := c.Apply(`{"score":3.14}`)
	if got != `{"score":"3.14"}` {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestCast_Apply_StringToBool(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "active", Type: "bool"}})
	got := c.Apply(`{"active":"true"}`)
	if got != `{"active":true}` {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestCast_Apply_StringToFloat(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "lat", Type: "float"}})
	got := c.Apply(`{"lat":"51.5"}`)
	if got != `{"lat":51.5}` {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestCast_Apply_MissingKeyNoOp(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "missing", Type: "int"}})
	input := `{"other":"value"}`
	got := c.Apply(input)
	if got == "" {
		t.Fatal("got empty string")
	}
}

func TestCast_Apply_CaseInsensitiveKey(t *testing.T) {
	c, _ := NewCast([]CastRule{{Key: "Count", Type: "int"}})
	got := c.Apply(`{"count":"5"}`)
	if got != `{"count":5}` {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestParseCastFlag_Valid(t *testing.T) {
	r, err := ParseCastFlag("duration:float")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Key != "duration" || r.Type != "float" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseCastFlag_Invalid(t *testing.T) {
	for _, s := range []string{"", "nocoion", ":int", "key:"} {
		_, err := ParseCastFlag(s)
		if err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}
