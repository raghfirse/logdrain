package formatter

import (
	"testing"
)

func TestNewReplace_NoRules(t *testing.T) {
	_, err := NewReplace(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestReplace_Apply_NonJSON(t *testing.T) {
	rule, _ := ParseReplaceFlag("msg/foo/bar")
	r, _ := NewReplace([]ReplaceRule{rule})
	line := "plain text"
	if got := r.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestReplace_Apply_ReplacesMatchingField(t *testing.T) {
	rule, _ := ParseReplaceFlag("msg/hello/world")
	r, _ := NewReplace([]ReplaceRule{rule})
	got := r.Apply(`{"msg":"say hello please"}`)
	want := `{"msg":"say world please"}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplace_Apply_CaseInsensitiveKey(t *testing.T) {
	rule, _ := ParseReplaceFlag("MSG/hello/world")
	r, _ := NewReplace([]ReplaceRule{rule})
	got := r.Apply(`{"msg":"hello"}`)
	want := `{"msg":"world"}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplace_Apply_NoMatchingField(t *testing.T) {
	rule, _ := ParseReplaceFlag("other/foo/bar")
	r, _ := NewReplace([]ReplaceRule{rule})
	line := `{"msg":"foo"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestReplace_Apply_NonStringFieldSkipped(t *testing.T) {
	rule, _ := ParseReplaceFlag("count/1/2")
	r, _ := NewReplace([]ReplaceRule{rule})
	line := `{"count":42}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestParsReplaceFlag_Valid(t *testing.T) {
	rule, err := ParseReplaceFlag("field/pat.*/repl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Key != "field" {
		t.Errorf("key: got %q, want %q", rule.Key, "field")
	}
	if rule.With != "repl" {
		t.Errorf("with: got %q, want %q", rule.With, "repl")
	}
}

func TestParseReplaceFlag_InvalidFormat(t *testing.T) {
	_, err := ParseReplaceFlag("noslash")
	if err == nil {
		t.Fatal("expected error for missing slashes")
	}
}

func TestParseReplaceFlag_EmptyKey(t *testing.T) {
	_, err := ParseReplaceFlag("/pat/repl")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParseReplaceFlag_InvalidPattern(t *testing.T) {
	_, err := ParseReplaceFlag("key/[invalid/repl")
	if err == nil {
		t.Fatal("expected error for bad regex")
	}
}

func TestReplaceFlagValue_SetAndString(t *testing.T) {
	f := NewReplaceFlagValue()
	if err := f.Set("msg/foo/bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules()) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(f.Rules()))
	}
	if f.String() == "" {
		t.Error("expected non-empty string")
	}
}

func TestReplaceFlagValue_Type(t *testing.T) {
	f := NewReplaceFlagValue()
	if f.Type() == "" {
		t.Error("expected non-empty type")
	}
}
