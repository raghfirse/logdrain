package formatter

import (
	"testing"
)

func TestExistsFlagValue_Set_Valid(t *testing.T) {
	f := NewExistsFlagValue()
	if err := f.Set("user,host"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg := f.Config()
	if len(cfg.Required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(cfg.Required))
	}
}

func TestExistsFlagValue_Set_Forbidden(t *testing.T) {
	f := NewExistsFlagValue()
	if err := f.Set("!debug"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg := f.Config()
	if len(cfg.Forbidden) != 1 || cfg.Forbidden[0] != "debug" {
		t.Errorf("unexpected forbidden: %v", cfg.Forbidden)
	}
}

func TestExistsFlagValue_Set_Empty(t *testing.T) {
	f := NewExistsFlagValue()
	if err := f.Set(""); err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestExistsFlagValue_Set_Whitespace(t *testing.T) {
	f := NewExistsFlagValue()
	if err := f.Set("   "); err == nil {
		t.Fatal("expected error for whitespace-only value")
	}
}

func TestExistsFlagValue_String_Default(t *testing.T) {
	f := NewExistsFlagValue()
	if got := f.String(); got != "" {
		t.Errorf("expected empty string default, got %q", got)
	}
}

func TestExistsFlagValue_String_AfterSet(t *testing.T) {
	f := NewExistsFlagValue()
	_ = f.Set("user,!debug")
	if got := f.String(); got != "user,!debug" {
		t.Errorf("expected 'user,!debug', got %q", got)
	}
}

func TestExistsFlagValue_Type(t *testing.T) {
	f := NewExistsFlagValue()
	if got := f.Type(); got != "exists" {
		t.Errorf("expected type 'exists', got %q", got)
	}
}
