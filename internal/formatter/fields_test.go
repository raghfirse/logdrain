package formatter

import (
	"testing"
)

func TestParseFieldsFlag_Empty(t *testing.T) {
	if got := ParseFieldsFlag(""); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestParseFieldsFlag_Single(t *testing.T) {
	got := ParseFieldsFlag("level")
	if len(got) != 1 || got[0] != "level" {
		t.Fatalf("unexpected %v", got)
	}
}

func TestParseFieldsFlag_Multiple(t *testing.T) {
	got := ParseFieldsFlag("level, msg , time")
	if len(got) != 3 {
		t.Fatalf("expected 3 fields, got %v", got)
	}
	if got[1] != "msg" {
		t.Fatalf("expected msg, got %s", got[1])
	}
}

func TestFieldFilter_AllowNoRules(t *testing.T) {
	ff := NewFieldFilter(nil, nil)
	if !ff.Allow("anything") {
		t.Fatal("expected allow with no rules")
	}
	if !ff.IsEmpty() {
		t.Fatal("expected IsEmpty true")
	}
}

func TestFieldFilter_IncludeOnly(t *testing.T) {
	ff := NewFieldFilter([]string{"level", "msg"}, nil)
	if !ff.Allow("level") {
		t.Fatal("level should be allowed")
	}
	if ff.Allow("time") {
		t.Fatal("time should not be allowed")
	}
}

func TestFieldFilter_ExcludeTakesPrecedence(t *testing.T) {
	ff := NewFieldFilter([]string{"level"}, []string{"level"})
	if ff.Allow("level") {
		t.Fatal("excluded field should not be allowed even if included")
	}
}

func TestFieldFilter_CaseInsensitive(t *testing.T) {
	ff := NewFieldFilter([]string{"Level"}, nil)
	if !ff.Allow("LEVEL") {
		t.Fatal("field matching should be case-insensitive")
	}
}

func TestFieldFilter_ExcludeOnly(t *testing.T) {
	ff := NewFieldFilter(nil, []string{"debug"})
	if ff.Allow("debug") {
		t.Fatal("excluded field should not be allowed")
	}
	if !ff.Allow("info") {
		t.Fatal("non-excluded field should be allowed")
	}
}
