package formatter

import (
	"testing"
)

func TestScriptFlag_Set_Valid(t *testing.T) {
	var f ScriptFlag
	if err := f.Set(`env=production`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules()) != 1 {
		t.Errorf("expected 1 rule, got %d", len(f.Rules()))
	}
}

func TestScriptFlag_Set_Multiple(t *testing.T) {
	var f ScriptFlag
	f.Set(`env=production`)
	f.Set(`app={{.service}}`)
	if len(f.Rules()) != 2 {
		t.Errorf("expected 2 rules, got %d", len(f.Rules()))
	}
}

func TestScriptFlag_Set_Empty(t *testing.T) {
	var f ScriptFlag
	if err := f.Set(""); err == nil {
		t.Fatal("expected error for empty rule")
	}
}

func TestScriptFlag_Set_Invalid(t *testing.T) {
	var f ScriptFlag
	if err := f.Set(`noequalssign`); err == nil {
		t.Fatal("expected error for rule without '='")
	}
}

func TestScriptFlag_String_Default(t *testing.T) {
	var f ScriptFlag
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestScriptFlag_String_AfterSet(t *testing.T) {
	var f ScriptFlag
	f.Set(`env=prod`)
	if f.String() == "" {
		t.Error("expected non-empty string after Set")
	}
}

func TestScriptFlag_Type(t *testing.T) {
	var f ScriptFlag
	if f.Type() != "key=template" {
		t.Errorf("unexpected type: %q", f.Type())
	}
}

func TestParseScriptFlag_Nil(t *testing.T) {
	s, err := ParseScriptFlag(nil)
	if err != nil || s != nil {
		t.Errorf("expected nil script and nil error, got %v, %v", s, err)
	}
}

func TestParseScriptFlag_Valid(t *testing.T) {
	var f ScriptFlag
	f.Set(`env=production`)
	s, err := ParseScriptFlag(&f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Script")
	}
}
