package formatter

import (
	"strings"
	"testing"
)

func TestWrapper_Disabled(t *testing.T) {
	w := NewWrapper(0, "  ")
	input := strings.Repeat("a", 200)
	if got := w.Wrap(input); got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestWrapper_ShortValue(t *testing.T) {
	w := NewWrapper(80, "  ")
	input := "short"
	if got := w.Wrap(input); got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestWrapper_LongValue_BreaksAtWidth(t *testing.T) {
	w := NewWrapper(10, "")
	input := "abcdefghijklmnopqrstu"
	got := w.Wrap(input)
	lines := strings.Split(got, "\n")
	for _, l := range lines {
		if len(l) > 10 {
			t.Errorf("line too long (%d): %q", len(l), l)
		}
	}
}

func TestWrapper_BreaksOnSpace(t *testing.T) {
	w := NewWrapper(10, "")
	input := "hello world foo"
	got := w.Wrap(input)
	if !strings.Contains(got, "\n") {
		t.Errorf("expected newline in wrapped output, got %q", got)
	}
	if strings.HasPrefix(got, "\n") {
		t.Errorf("unexpected leading newline: %q", got)
	}
}

func TestWrapper_IndentApplied(t *testing.T) {
	w := NewWrapper(5, ">>")
	input := "abcdefghij"
	got := w.Wrap(input)
	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %q", got)
	}
	for _, l := range lines[1:] {
		if !strings.HasPrefix(l, ">>") {
			t.Errorf("expected indent prefix, got %q", l)
		}
	}
}

func TestParseWrapFlag_Valid(t *testing.T) {
	cases := []struct{ in string; want int }{
		{"", 0}, {"0", 0}, {"off", 0}, {"80", 80}, {"120", 120},
	}
	for _, c := range cases {
		got, err := ParseWrapFlag(c.in)
		if err != nil {
			t.Errorf("ParseWrapFlag(%q) error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("ParseWrapFlag(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseWrapFlag_Invalid(t *testing.T) {
	for _, s := range []string{"abc", "-1", "1.5"} {
		if _, err := ParseWrapFlag(s); err == nil {
			t.Errorf("ParseWrapFlag(%q) expected error", s)
		}
	}
}
