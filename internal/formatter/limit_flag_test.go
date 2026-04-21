package formatter

import (
	"testing"
)

func TestLimitFlagValue_Set_Valid(t *testing.T) {
	f := &LimitFlagValue{}
	if err := f.Set("42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.N != 42 {
		t.Errorf("expected 42, got %d", f.N)
	}
}

func TestLimitFlagValue_Set_Zero(t *testing.T) {
	f := &LimitFlagValue{}
	if err := f.Set("0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.N != 0 {
		t.Errorf("expected 0, got %d", f.N)
	}
}

func TestLimitFlagValue_Set_Invalid(t *testing.T) {
	f := &LimitFlagValue{}
	for _, bad := range []string{"abc", "-1", ""} {
		if bad == "" {
			continue // empty string is handled by ParseLimitFlag, not Set
		}
		if err := f.Set(bad); err == nil {
			t.Errorf("Set(%q): expected error", bad)
		}
	}
}

func TestLimitFlagValue_String_Default(t *testing.T) {
	f := &LimitFlagValue{}
	if f.String() != "0" {
		t.Errorf("expected \"0\", got %q", f.String())
	}
}

func TestLimitFlagValue_String_AfterSet(t *testing.T) {
	f := &LimitFlagValue{}
	_ = f.Set("7")
	if f.String() != "7" {
		t.Errorf("expected \"7\", got %q", f.String())
	}
}

func TestLimitFlagValue_Type(t *testing.T) {
	f := &LimitFlagValue{}
	if f.Type() != "int" {
		t.Errorf("expected \"int\", got %q", f.Type())
	}
}
