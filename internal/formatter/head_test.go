package formatter

import (
	"testing"
)

func TestHead_ZeroLimitAllowsAll(t *testing.T) {
	h, err := NewHead(0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !h.Allow(`{"msg":"hi"}`) {
			t.Fatalf("expected Allow=true on iteration %d", i)
		}
	}
}

func TestHead_LimitGlobal(t *testing.T) {
	h, _ := NewHead(3, "")
	results := make([]bool, 5)
	for i := range results {
		results[i] = h.Allow(`{"msg":"line"}`)
	}
	expected := []bool{true, true, true, false, false}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, got, expected[i])
		}
	}
}

func TestHead_LimitPerKey(t *testing.T) {
	h, _ := NewHead(2, "svc")
	lines := []string{
		`{"svc":"a"}`,
		`{"svc":"a"}`,
		`{"svc":"a"}`,
		`{"svc":"b"}`,
		`{"svc":"b"}`,
	}
	expected := []bool{true, true, false, true, true}
	for i, line := range lines {
		if got := h.Allow(line); got != expected[i] {
			t.Errorf("line %d: got %v, want %v", i, got, expected[i])
		}
	}
}

func TestHead_Reset(t *testing.T) {
	h, _ := NewHead(1, "")
	if !h.Allow(`{}`) {
		t.Fatal("first Allow should be true")
	}
	if h.Allow(`{}`) {
		t.Fatal("second Allow should be false")
	}
	h.Reset()
	if !h.Allow(`{}`) {
		t.Fatal("Allow after Reset should be true")
	}
}

func TestHead_NegativeLimitError(t *testing.T) {
	_, err := NewHead(-1, "")
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestParseHeadFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		limit int
		key   string
	}{
		{"10", 10, ""},
		{"5:service", 5, "service"},
		{"0", 0, ""},
	}
	for _, tc := range cases {
		h, err := ParseHeadFlag(tc.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", tc.input, err)
			continue
		}
		if h.limit != tc.limit {
			t.Errorf("%q: limit got %d, want %d", tc.input, h.limit, tc.limit)
		}
		if h.key != tc.key {
			t.Errorf("%q: key got %q, want %q", tc.input, h.key, tc.key)
		}
	}
}

func TestParseHeadFlag_Invalid(t *testing.T) {
	cases := []string{"", "abc", "abc:svc", "-3"}
	for _, tc := range cases {
		_, err := ParseHeadFlag(tc)
		if err == nil {
			t.Errorf("%q: expected error", tc)
		}
	}
}

func TestHeadFlagValue_SetAndString(t *testing.T) {
	var f HeadFlagValue
	if err := f.Set("4:env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.String() != "4:env" {
		t.Errorf("String() got %q, want %q", f.String(), "4:env")
	}
	if f.Head() == nil {
		t.Error("Head() should not be nil after Set")
	}
}

func TestHeadFlagValue_Type(t *testing.T) {
	var f HeadFlagValue
	if f.Type() != "head" {
		t.Errorf("Type() got %q, want \"head\"", f.Type())
	}
}
