package formatter

import (
	"testing"
)

func TestFlattenFlag_Set_Valid(t *testing.T) {
	var f FlattenFlag
	if err := f.Set("."); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsSet() {
		t.Error("expected IsSet() to be true")
	}
	if f.Value == nil {
		t.Error("expected Value to be non-nil")
	}
}

func TestFlattenFlag_Set_Empty(t *testing.T) {
	var f FlattenFlag
	if err := f.Set(""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value.separator != "." {
		t.Errorf("expected default separator '.', got %q", f.Value.separator)
	}
}

func TestFlattenFlag_Set_Invalid(t *testing.T) {
	var f FlattenFlag
	if err := f.Set("\t"); err == nil {
		t.Error("expected error for tab separator")
	}
}

func TestFlattenFlag_String_Default(t *testing.T) {
	var f FlattenFlag
	if got := f.String(); got != "" {
		t.Errorf("expected empty string before set, got %q", got)
	}
}

func TestFlattenFlag_String_AfterSet(t *testing.T) {
	var f FlattenFlag
	_ = f.Set("/")
	got := f.String()
	if got == "" {
		t.Error("expected non-empty string after set")
	}
}

func TestFlattenFlag_Type(t *testing.T) {
	var f FlattenFlag
	if got := f.Type(); got != "separator" {
		t.Errorf("expected 'separator', got %q", got)
	}
}
