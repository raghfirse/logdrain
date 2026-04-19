package formatter

import "testing"

func TestNormalizeLevel(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"INFO", "info"},
		{"warning", "warn"},
		{"WARNING", "warn"},
		{"err", "error"},
		{"critical", "fatal"},
		{"panic", "fatal"},
		{"debug", "debug"},
		{"trace", "trace"},
	}
	for _, tc := range cases {
		got := NormalizeLevel(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeLevel(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseLevelFilter_Valid(t *testing.T) {
	f, err := ParseLevelFilter("warn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f("debug") {
		t.Error("debug should not pass warn filter")
	}
	if f("info") {
		t.Error("info should not pass warn filter")
	}
	if !f("warn") {
		t.Error("warn should pass warn filter")
	}
	if !f("error") {
		t.Error("error should pass warn filter")
	}
	if !f("fatal") {
		t.Error("fatal should pass warn filter")
	}
}

func TestParseLevelFilter_UnknownLevelPassesThrough(t *testing.T) {
	f, err := ParseLevelFilter("info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f("someunknownlevel") {
		t.Error("unknown level should pass through")
	}
}

func TestParseLevelFilter_Invalid(t *testing.T) {
	_, err := ParseLevelFilter("verbose")
	if err == nil {
		t.Error("expected error for unknown min level")
	}
}

func TestParseLevelFilter_NormalizesInput(t *testing.T) {
	f, err := ParseLevelFilter("WARNING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f("info") {
		t.Error("info should not pass warn filter")
	}
	if !f("error") {
		t.Error("error should pass warn filter")
	}
}
