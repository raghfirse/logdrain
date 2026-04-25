package formatter

import (
	"testing"
)

func TestBetween_SuppressesBeforeStart(t *testing.T) {
	b, _ := NewBetween("START", "END", "")
	if got := b.Apply(`{"msg":"before"}`); got != "" {
		t.Fatalf("expected suppressed, got %q", got)
	}
}

func TestBetween_IncludesStartLine(t *testing.T) {
	b, _ := NewBetween("START", "END", "")
	line := "START marker"
	if got := b.Apply(line); got != line {
		t.Fatalf("expected start line passed through, got %q", got)
	}
}

func TestBetween_PassesThroughActiveLines(t *testing.T) {
	b, _ := NewBetween("START", "END", "")
	b.Apply("START")
	middle := "middle line"
	if got := b.Apply(middle); got != middle {
		t.Fatalf("expected middle line, got %q", got)
	}
}

func TestBetween_IncludesEndLine(t *testing.T) {
	b, _ := NewBetween("START", "END", "")
	b.Apply("START")
	end := "END marker"
	if got := b.Apply(end); got != end {
		t.Fatalf("expected end line passed through, got %q", got)
	}
}

func TestBetween_SuppressesAfterEnd(t *testing.T) {
	b, _ := NewBetween("START", "END", "")
	b.Apply("START")
	b.Apply("END")
	if got := b.Apply("after"); got != "" {
		t.Fatalf("expected suppressed after end, got %q", got)
	}
}

func TestBetween_FieldResolution(t *testing.T) {
	b, _ := NewBetween("begin", "finish", "event")
	if got := b.Apply(`{"event":"begin here"}`); got == "" {
		t.Fatal("expected start matched on field")
	}
	if got := b.Apply(`{"event":"finish now"}`); got == "" {
		t.Fatal("expected end line passed through")
	}
	if got := b.Apply(`{"event":"after"}`); got != "" {
		t.Fatalf("expected suppressed after end, got %q", got)
	}
}

func TestBetween_Reset(t *testing.T) {
	b, _ := NewBetween("START", "END", "")
	b.Apply("START")
	b.Reset()
	if got := b.Apply("middle"); got != "" {
		t.Fatalf("expected suppressed after reset, got %q", got)
	}
}

func TestBetween_EmptyStartError(t *testing.T) {
	_, err := NewBetween("", "END", "")
	if err == nil {
		t.Fatal("expected error for empty start pattern")
	}
}

func TestBetween_EmptyEndError(t *testing.T) {
	_, err := NewBetween("START", "", "")
	if err == nil {
		t.Fatal("expected error for empty end pattern")
	}
}

func TestParseBetweenFlag_Valid(t *testing.T) {
	b, err := ParseBetweenFlag("foo:bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.startPattern != "foo" || b.endPattern != "bar" {
		t.Fatalf("unexpected patterns: %+v", b)
	}
}

func TestParseBetweenFlag_WithField(t *testing.T) {
	b, err := ParseBetweenFlag("foo:bar:msg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.field != "msg" {
		t.Fatalf("expected field 'msg', got %q", b.field)
	}
}

func TestParseBetweenFlag_Invalid(t *testing.T) {
	_, err := ParseBetweenFlag("onlyone")
	if err == nil {
		t.Fatal("expected error for missing end pattern")
	}
}
