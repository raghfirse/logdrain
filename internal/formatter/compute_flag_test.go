package formatter

import (
	"testing"
)

func TestComputeFlag_Set_Valid(t *testing.T) {
	f := &computeFlagValue{}
	if err := f.Set("total=price+tax"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.exprs) != 1 {
		t.Fatalf("expected 1 expr, got %d", len(f.exprs))
	}
}

func TestComputeFlag_Set_Multiple(t *testing.T) {
	f := &computeFlagValue{}
	_ = f.Set("total=a+b")
	_ = f.Set("diff=a-b")
	if len(f.exprs) != 2 {
		t.Fatalf("expected 2 exprs, got %d", len(f.exprs))
	}
}

func TestComputeFlag_Set_Empty(t *testing.T) {
	f := &computeFlagValue{}
	if err := f.Set(""); err == nil {
		t.Fatal("expected error for empty expression")
	}
}

func TestComputeFlag_Set_Invalid(t *testing.T) {
	f := &computeFlagValue{}
	if err := f.Set("badexpr"); err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestComputeFlag_String_Default(t *testing.T) {
	f := &computeFlagValue{}
	if f.String() != "" {
		t.Fatalf("expected empty string, got %q", f.String())
	}
}

func TestComputeFlag_String_AfterSet(t *testing.T) {
	f := &computeFlagValue{}
	_ = f.Set("total=a+b")
	if f.String() == "" {
		t.Fatal("expected non-empty string after Set")
	}
}

func TestComputeFlag_Type(t *testing.T) {
	f := &computeFlagValue{}
	if f.Type() != "expr" {
		t.Fatalf("expected type 'expr', got %q", f.Type())
	}
}

func TestParseComputeFlag_Nil(t *testing.T) {
	c, err := ParseComputeFlag(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Compute")
	}
}

func TestParseComputeFlag_WithExprs(t *testing.T) {
	f := &computeFlagValue{}
	_ = f.Set("out=x*y")
	c, err := ParseComputeFlag(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := c.Apply(`{"x":3,"y":4}`)
	assertJSONField(t, result, "out", 12.0)
}
