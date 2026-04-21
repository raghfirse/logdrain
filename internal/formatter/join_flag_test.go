package formatter

import (
	"testing"
)

func TestJoinFlagValue_Set_Valid(t *testing.T) {
	var f JoinFlagValue
	if err := f.Set("first+last->full"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules()) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(f.Rules()))
	}
}

func TestJoinFlagValue_Set_Multiple(t *testing.T) {
	var f JoinFlagValue
	_ = f.Set("a+b->c")
	_ = f.Set("x+y->z")
	if len(f.Rules()) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(f.Rules()))
	}
}

func TestJoinFlagValue_Set_Invalid(t *testing.T) {
	var f JoinFlagValue
	if err := f.Set("bad_input"); err == nil {
		t.Fatal("expected error for invalid input")
	}
}

func TestJoinFlagValue_String_Default(t *testing.T) {
	var f JoinFlagValue
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestJoinFlagValue_String_AfterSet(t *testing.T) {
	var f JoinFlagValue
	_ = f.Set("host+port:separator->addr")
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string after Set")
	}
}

func TestJoinFlagValue_Type(t *testing.T) {
	var f JoinFlagValue
	if f.Type() != "join" {
		t.Errorf("expected 'join', got %q", f.Type())
	}
}
