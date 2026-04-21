package formatter

import (
	"testing"
)

func TestGrepFlag_Set_PatternOnly(t *testing.T) {
	var f GrepFlag
	if err := f.Set("error"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Grep() == nil {
		t.Fatal("expected non-nil Grep")
	}
	if f.String() != "error" {
		t.Errorf("got %q, want %q", f.String(), "error")
	}
}

func TestGrepFlag_Set_PatternWithFields(t *testing.T) {
	var f GrepFlag
	if err := f.Set("boom:msg,service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.String() != "boom:msg,service" {
		t.Errorf("got %q, want %q", f.String(), "boom:msg,service")
	}
	g := f.Grep()
	if g == nil {
		t.Fatal("expected non-nil Grep")
	}
	if !g.Match(`{"msg":"boom"}`) {
		t.Error("expected match")
	}
}

func TestGrepFlag_Set_Invalid(t *testing.T) {
	var f GrepFlag
	if err := f.Set("[bad"); err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestGrepFlag_String_Default(t *testing.T) {
	var f GrepFlag
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestGrepFlag_Type(t *testing.T) {
	var f GrepFlag
	if f.Type() != "grep" {
		t.Errorf("expected 'grep', got %q", f.Type())
	}
}

func TestParseGrepFlag_Valid(t *testing.T) {
	f, err := ParseGrepFlag("warn:level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Grep() == nil {
		t.Fatal("expected non-nil Grep")
	}
}

func TestParseGrepFlag_Empty(t *testing.T) {
	f, err := ParseGrepFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Grep() != nil {
		t.Error("expected nil Grep for empty flag")
	}
}

func TestParseGrepFlag_Invalid(t *testing.T) {
	_, err := ParseGrepFlag("[nope")
	if err == nil {
		t.Fatal("expected error")
	}
}
