package formatter

import "testing"

func TestLevelFlag_Set_Valid(t *testing.T) {
	f := &LevelFlag{}
	if err := f.Set("warn"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != "warn" {
		t.Errorf("expected warn, got %q", f.Value)
	}
}

func TestLevelFlag_Set_Normalizes(t *testing.T) {
	f := &LevelFlag{}
	if err := f.Set("WARNING"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value != "warn" {
		t.Errorf("expected warn, got %q", f.Value)
	}
}

func TestLevelFlag_Set_Invalid(t *testing.T) {
	f := &LevelFlag{}
	if err := f.Set("verbose"); err == nil {
		t.Error("expected error for invalid level")
	}
}

func TestLevelFlag_String_Default(t *testing.T) {
	f := &LevelFlag{}
	if f.String() != "info" {
		t.Errorf("expected default info, got %q", f.String())
	}
}

func TestLevelFlag_Type(t *testing.T) {
	f := &LevelFlag{}
	if f.Type() != "level" {
		t.Errorf("expected level, got %q", f.Type())
	}
}
