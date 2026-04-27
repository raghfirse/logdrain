package formatter

import (
	"testing"
)

func TestOmitFlagValue_Set_Valid(t *testing.T) {
	f := NewOmitFlagValue()
	if err := f.Set("foo,bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Keys()) != 2 {
		t.Errorf("expected 2 keys, got %d", len(f.Keys()))
	}
}

func TestOmitFlagValue_Set_Multiple(t *testing.T) {
	f := NewOmitFlagValue()
	_ = f.Set("a")
	_ = f.Set("b,c")
	if len(f.Keys()) != 3 {
		t.Errorf("expected 3 accumulated keys, got %d", len(f.Keys()))
	}
}

func TestOmitFlagValue_Set_Empty(t *testing.T) {
	f := NewOmitFlagValue()
	if err := f.Set(""); err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestOmitFlagValue_String_Default(t *testing.T) {
	f := NewOmitFlagValue()
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestOmitFlagValue_String_AfterSet(t *testing.T) {
	f := NewOmitFlagValue()
	_ = f.Set("x,y")
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string after Set")
	}
}

func TestOmitFlagValue_Type(t *testing.T) {
	f := NewOmitFlagValue()
	if f.Type() == "" {
		t.Error("expected non-empty type string")
	}
}
