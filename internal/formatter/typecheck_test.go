package formatter

import (
	"testing"
)

func TestNewTypeChecker_Valid(t *testing.T) {
	for _, typ := range []string{"string", "number", "bool", "array", "object", "null"} {
		_, err := NewTypeChecker("field", typ, false)
		if err != nil {
			t.Errorf("unexpected error for type %q: %v", typ, err)
		}
	}
}

func TestNewTypeChecker_EmptyKey(t *testing.T) {
	_, err := NewTypeChecker("", "string", false)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewTypeChecker_UnknownType(t *testing.T) {
	_, err := NewTypeChecker("k", "integer", false)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestTypeChecker_Apply_MatchesString(t *testing.T) {
	c, _ := NewTypeChecker("msg", "string", false)
	got := c.Apply(`{"msg":"hello"}`)
	if got == "" {
		t.Fatal("expected line to pass")
	}
}

func TestTypeChecker_Apply_WrongType(t *testing.T) {
	c, _ := NewTypeChecker("msg", "number", false)
	got := c.Apply(`{"msg":"hello"}`)
	if got != "" {
		t.Fatalf("expected line to be suppressed, got %q", got)
	}
}

func TestTypeChecker_Apply_Inverted(t *testing.T) {
	c, _ := NewTypeChecker("count", "string", true)
	got := c.Apply(`{"count":42}`)
	if got == "" {
		t.Fatal("expected inverted match to pass")
	}
}

func TestTypeChecker_Apply_NonJSONPassthrough(t *testing.T) {
	c, _ := NewTypeChecker("k", "bool", false)
	line := "not json at all"
	if got := c.Apply(line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestTypeChecker_Apply_MissingFieldSuppressed(t *testing.T) {
	c, _ := NewTypeChecker("missing", "string", false)
	got := c.Apply(`{"other":"val"}`)
	if got != "" {
		t.Fatalf("expected suppression for missing field, got %q", got)
	}
}

func TestTypeChecker_Apply_MissingFieldInverted(t *testing.T) {
	c, _ := NewTypeChecker("missing", "string", true)
	got := c.Apply(`{"other":"val"}`)
	if got == "" {
		t.Fatal("expected inverted missing field to pass")
	}
}

func TestParseTypeCheckFlag_Valid(t *testing.T) {
	c, err := ParseTypeCheckFlag("level:string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.key != "level" || c.wantType != "string" || c.invert {
		t.Fatalf("unexpected checker: %+v", c)
	}
}

func TestParseTypeCheckFlag_Inverted(t *testing.T) {
	c, err := ParseTypeCheckFlag("!count:number")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.invert {
		t.Fatal("expected invert=true")
	}
}

func TestParseTypeCheckFlag_Invalid(t *testing.T) {
	if _, err := ParseTypeCheckFlag("nocolon"); err == nil {
		t.Fatal("expected error for missing colon")
	}
	if _, err := ParseTypeCheckFlag("k:badtype"); err == nil {
		t.Fatal("expected error for bad type")
	}
}

func TestTypeCheckFlagValue_SetAndString(t *testing.T) {
	f := NewTypeCheckFlagValue()
	if f.String() != "none" {
		t.Fatalf("expected 'none', got %q", f.String())
	}
	if err := f.Set("level:string"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Checkers()) != 1 {
		t.Fatalf("expected 1 checker, got %d", len(f.Checkers()))
	}
	if f.String() != "level:string" {
		t.Fatalf("unexpected string: %q", f.String())
	}
}
