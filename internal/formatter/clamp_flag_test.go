package formatter

import (
	"testing"
)

func TestClampFlagValue_Set_Valid(t *testing.T) {
	f := NewClampFlagValue()
	if err := f.Set("score,0,100"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Clamps()) != 1 {
		t.Fatalf("expected 1 clamp, got %d", len(f.Clamps()))
	}
}

func TestClampFlagValue_Set_Multiple(t *testing.T) {
	f := NewClampFlagValue()
	_ = f.Set("score,0,100")
	_ = f.Set("latency,0,5000")
	if len(f.Clamps()) != 2 {
		t.Fatalf("expected 2 clamps, got %d", len(f.Clamps()))
	}
}

func TestClampFlagValue_Set_Invalid(t *testing.T) {
	f := NewClampFlagValue()
	if err := f.Set("badspec"); err == nil {
		t.Fatal("expected error for bad spec")
	}
	if len(f.Clamps()) != 0 {
		t.Fatal("expected no clamps after error")
	}
}

func TestClampFlagValue_String_Default(t *testing.T) {
	f := NewClampFlagValue()
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestClampFlagValue_String_AfterSet(t *testing.T) {
	f := NewClampFlagValue()
	_ = f.Set("score,0,100")
	if f.String() != "score,0,100" {
		t.Errorf("unexpected string: %q", f.String())
	}
}

func TestClampFlagValue_Type(t *testing.T) {
	f := NewClampFlagValue()
	if f.Type() != "field,min,max" {
		t.Errorf("unexpected type: %q", f.Type())
	}
}
