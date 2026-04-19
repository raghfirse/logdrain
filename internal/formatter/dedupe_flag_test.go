package formatter

import (
	"testing"
	"time"
)

func TestParseDedupeFlag_Valid(t *testing.T) {
	f, err := ParseDedupeFlag("5s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Duration != 5*time.Second {
		t.Fatalf("expected 5s, got %v", f.Duration)
	}
	if !f.Set {
		t.Fatal("expected Set=true")
	}
}

func TestParseDedupeFlag_Empty(t *testing.T) {
	f, err := ParseDedupeFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Set {
		t.Fatal("expected Set=false for empty string")
	}
}

func TestParseDedupeFlag_Zero(t *testing.T) {
	f, err := ParseDedupeFlag("0s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Set {
		t.Fatal("expected Set=false for zero duration")
	}
}

func TestParseDedupeFlag_Invalid(t *testing.T) {
	_, err := ParseDedupeFlag("notaduration")
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestParseDedupeFlag_Negative(t *testing.T) {
	_, err := ParseDedupeFlag("-1s")
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestDedupeFlag_Value_ReturnsNilWhenNotSet(t *testing.T) {
	f := DedupeFlag{}
	if f.Value() != nil {
		t.Fatal("expected nil filter when not set")
	}
}

func TestDedupeFlag_Value_ReturnsFilter(t *testing.T) {
	f := DedupeFlag{Duration: 2 * time.Second, Set: true}
	if f.Value() == nil {
		t.Fatal("expected non-nil filter")
	}
}
