package formatter

import (
	"testing"
)

func TestTruncator_Disabled(t *testing.T) {
	tr := NewTruncator(0, "")
	fields := map[string]any{"msg": "hello world"}
	out := tr.Apply(fields)
	if out["msg"] != "hello world" {
		t.Fatalf("expected unchanged, got %v", out["msg"])
	}
}

func TestTruncator_ShortValue(t *testing.T) {
	tr := NewTruncator(20, "…")
	fields := map[string]any{"msg": "hi"}
	out := tr.Apply(fields)
	if out["msg"] != "hi" {
		t.Fatalf("expected unchanged, got %v", out["msg"])
	}
}

func TestTruncator_LongValue(t *testing.T) {
	tr := NewTruncator(5, "…")
	fields := map[string]any{"msg": "hello world"}
	out := tr.Apply(fields)
	if out["msg"] != "hello…" {
		t.Fatalf("expected truncated, got %v", out["msg"])
	}
}

func TestTruncator_NonStringUnchanged(t *testing.T) {
	tr := NewTruncator(5, "…")
	fields := map[string]any{"count": 42}
	out := tr.Apply(fields)
	if out["count"] != 42 {
		t.Fatalf("expected 42, got %v", out["count"])
	}
}

func TestTruncator_CustomSuffix(t *testing.T) {
	tr := NewTruncator(3, "...")
	fields := map[string]any{"k": "abcdef"}
	out := tr.Apply(fields)
	if out["k"] != "abc..." {
		t.Fatalf("expected abc..., got %v", out["k"])
	}
}

func TestParseTruncateFlag_Valid(t *testing.T) {
	maxLen, suffix, err := ParseTruncateFlag("80")
	if err != nil || maxLen != 80 || suffix != "…" {
		t.Fatalf("unexpected: %d %q %v", maxLen, suffix, err)
	}
}

func TestParseTruncateFlag_CustomSuffix(t *testing.T) {
	maxLen, suffix, err := ParseTruncateFlag("10:>>")
	if err != nil || maxLen != 10 || suffix != ">>" {
		t.Fatalf("unexpected: %d %q %v", maxLen, suffix, err)
	}
}

func TestParseTruncateFlag_Zero(t *testing.T) {
	maxLen, _, err := ParseTruncateFlag("0")
	if err != nil || maxLen != 0 {
		t.Fatalf("unexpected: %d %v", maxLen, err)
	}
}

func TestParseTruncateFlag_Invalid(t *testing.T) {
	_, _, err := ParseTruncateFlag("abc")
	if err == nil {
		t.Fatal("expected error")
	}
}
