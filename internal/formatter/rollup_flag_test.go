package formatter

import (
	"testing"
)

func TestRollupFlagValue_Set_Valid(t *testing.T) {
	var f RollupFlagValue
	if err := f.Set("service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "service" {
		t.Errorf("expected service, got %q", f.Field())
	}
	if !f.IsSet() {
		t.Error("expected IsSet to be true")
	}
}

func TestRollupFlagValue_Set_Empty(t *testing.T) {
	var f RollupFlagValue
	if err := f.Set(""); err == nil {
		t.Fatal("expected error for empty field")
	}
	if f.IsSet() {
		t.Error("expected IsSet to remain false after error")
	}
}

func TestRollupFlagValue_String_Default(t *testing.T) {
	var f RollupFlagValue
	if f.String() != "" {
		t.Errorf("expected empty string default, got %q", f.String())
	}
}

func TestRollupFlagValue_String_AfterSet(t *testing.T) {
	var f RollupFlagValue
	f.Set("region")
	if f.String() != "region" {
		t.Errorf("expected region, got %q", f.String())
	}
}

func TestRollupFlagValue_Type(t *testing.T) {
	var f RollupFlagValue
	if f.Type() != "field" {
		t.Errorf("expected 'field', got %q", f.Type())
	}
}

func TestRollupFlagUsage_NotEmpty(t *testing.T) {
	if RollupFlagUsage == "" {
		t.Error("expected non-empty usage string")
	}
}
