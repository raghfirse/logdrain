package formatter

import (
	"testing"
)

func TestLimiter_ZeroAllowsAll(t *testing.T) {
	l, err := NewLimiter(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 1000; i++ {
		out, done := l.Apply(`{"msg":"hi"}`)
		if done {
			t.Fatalf("expected done=false at i=%d", i)
		}
		if out == "" {
			t.Fatalf("expected non-empty output at i=%d", i)
		}
	}
}

func TestLimiter_StopsAtMax(t *testing.T) {
	l, _ := NewLimiter(3)
	line := `{"msg":"x"}`
	for i := 0; i < 3; i++ {
		out, done := l.Apply(line)
		if done {
			t.Fatalf("expected done=false at i=%d", i)
		}
		if out != line {
			t.Fatalf("expected line unchanged at i=%d", i)
		}
	}
	_, done := l.Apply(line)
	if !done {
		t.Fatal("expected done=true after limit reached")
	}
}

func TestLimiter_Reset(t *testing.T) {
	l, _ := NewLimiter(1)
	l.Apply(`{"a":1}`)
	_, done := l.Apply(`{"b":2}`)
	if !done {
		t.Fatal("expected done=true")
	}
	l.Reset()
	out, done := l.Apply(`{"c":3}`)
	if done {
		t.Fatal("expected done=false after reset")
	}
	if out == "" {
		t.Fatal("expected output after reset")
	}
}

func TestNewLimiter_NegativeError(t *testing.T) {
	_, err := NewLimiter(-1)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestParseLimitFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"0", 0},
		{"10", 10},
		{"1", 1},
	}
	for _, tc := range cases {
		got, err := ParseLimitFlag(tc.input)
		if err != nil {
			t.Errorf("ParseLimitFlag(%q) error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseLimitFlag(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestParseLimitFlag_Invalid(t *testing.T) {
	for _, s := range []string{"abc", "-5", "1.5"} {
		_, err := ParseLimitFlag(s)
		if err == nil {
			t.Errorf("ParseLimitFlag(%q): expected error", s)
		}
	}
}

// TestLimiter_ResetMultipleTimes verifies that a limiter can be reset and reused
// repeatedly, each time allowing exactly max lines through.
func TestLimiter_ResetMultipleTimes(t *testing.T) {
	const max = 2
	l, _ := NewLimiter(max)
	line := `{"msg":"test"}`

	for cycle := 0; cycle < 3; cycle++ {
		for i := 0; i < max; i++ {
			_, done := l.Apply(line)
			if done {
				t.Fatalf("cycle %d: expected done=false at i=%d", cycle, i)
			}
		}
		_, done := l.Apply(line)
		if !done {
			t.Fatalf("cycle %d: expected done=true after limit", cycle)
		}
		l.Reset()
	}
}
