package formatter

import (
	"testing"
)

func TestTopFlagValue_Set_Valid(t *testing.T) {
	var f TopFlagValue
	if err := f.Set("level:5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "level" {
		t.Errorf("expected field 'level', got %q", f.Field)
	}
	if f.N != 5 {
		t.Errorf("expected n=5, got %d", f.N)
	}
}

func TestTopFlagValue_Set_InvalidFormat(t *testing.T) {
	var f TopFlagValue
	if err := f.Set("level"); err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestTopFlagValue_Set_InvalidN(t *testing.T) {
	var f TopFlagValue
	if err := f.Set("level:0"); err == nil {
		t.Error("expected error for n=0")
	}
}

func TestTopFlagValue_Set_NegativeN(t *testing.T) {
	var f TopFlagValue
	if err := f.Set("level:-1"); err == nil {
		t.Error("expected error for negative n")
	}
}

func TestTopFlagValue_String_Default(t *testing.T) {
	var f TopFlagValue
	if f.String() != "" {
		t.Errorf("expected empty string by default, got %q", f.String())
	}
}

func TestTopFlagValue_String_AfterSet(t *testing.T) {
	var f TopFlagValue
	_ = f.Set("status:3")
	if f.String() != "status:3" {
		t.Errorf("expected 'status:3', got %q", f.String())
	}
}

func TestTopFlagValue_Type(t *testing.T) {
	var f TopFlagValue
	if f.Type() != "field:n" {
		t.Errorf("expected type 'field:n', got %q", f.Type())
	}
}

func TestTopFlagValue_IsSet(t *testing.T) {
	var f TopFlagValue
	if f.IsSet() {
		t.Error("expected IsSet=false before Set")
	}
	_ = f.Set("level:10")
	if !f.IsSet() {
		t.Error("expected IsSet=true after Set")
	}
}
