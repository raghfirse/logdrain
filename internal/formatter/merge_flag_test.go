package formatter

import (
	"testing"
)

func TestMergeFlagValue_Set_Valid(t *testing.T) {
	f := NewMergeFlagValue(false)
	if err := f.Set(`{"env":"prod"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Merger() == nil {
		t.Fatal("expected non-nil merger after Set")
	}
}

func TestMergeFlagValue_Set_Invalid(t *testing.T) {
	f := NewMergeFlagValue(false)
	if err := f.Set("not-json"); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestMergeFlagValue_String_Default(t *testing.T) {
	f := NewMergeFlagValue(false)
	if f.String() != "" {
		t.Errorf("expected empty string by default, got %q", f.String())
	}
}

func TestMergeFlagValue_String_AfterSet(t *testing.T) {
	f := NewMergeFlagValue(false)
	raw := `{"region":"eu-west-1"}`
	f.Set(raw)
	if f.String() != raw {
		t.Errorf("expected %q, got %q", raw, f.String())
	}
}

func TestMergeFlagValue_Type(t *testing.T) {
	f := NewMergeFlagValue(false)
	if f.Type() != "json-object" {
		t.Errorf("expected type 'json-object', got %q", f.Type())
	}
}

func TestMergeFlagValue_Overwrite(t *testing.T) {
	f := NewMergeFlagValue(true)
	f.Set(`{"env":"prod"}`)
	if f.merger == nil {
		t.Fatal("expected merger to be set")
	}
	if !f.overwrite {
		t.Error("expected overwrite=true")
	}
}

func TestMergeFlagUsage_NoOverwrite(t *testing.T) {
	s := MergeFlagUsage(false)
	if s == "" {
		t.Error("expected non-empty usage string")
	}
}

func TestMergeFlagUsage_Overwrite(t *testing.T) {
	s := MergeFlagUsage(true)
	if s == "" {
		t.Error("expected non-empty usage string")
	}
}
