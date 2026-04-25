package formatter

import (
	"testing"
)

func TestPickFlagValue_Set_Valid(t *testing.T) {
	f := NewPickFlagValue()
	if err := f.Set("msg,level"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := f.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "msg" || keys[1] != "level" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestPickFlagValue_Set_Multiple(t *testing.T) {
	f := NewPickFlagValue()
	_ = f.Set("msg")
	_ = f.Set("level")
	if len(f.Keys()) != 2 {
		t.Errorf("expected 2 accumulated keys, got %d", len(f.Keys()))
	}
}

func TestPickFlagValue_Set_Empty(t *testing.T) {
	f := NewPickFlagValue()
	if err := f.Set(""); err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestPickFlagValue_String_Default(t *testing.T) {
	f := NewPickFlagValue()
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestPickFlagValue_String_AfterSet(t *testing.T) {
	f := NewPickFlagValue()
	_ = f.Set("msg,level")
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string after Set")
	}
}

func TestPickFlagValue_Type(t *testing.T) {
	f := NewPickFlagValue()
	if f.Type() != "fields" {
		t.Errorf("expected 'fields', got %q", f.Type())
	}
}
