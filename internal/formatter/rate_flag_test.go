package formatter

import (
	"testing"
	"time"
)

func TestParseRateFlag_Valid(t *testing.T) {
	cases := []struct {
		input  string
		max    int
		win    time.Duration
	}{
		{"10/1s", 10, time.Second},
		{"5/500ms", 5, 500 * time.Millisecond},
		{"0/1m", 0, time.Minute},
		{"100/2m", 100, 2 * time.Minute},
	}
	for _, c := range cases {
		rl, err := ParseRateFlag(c.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.input, err)
			continue
		}
		if rl.Max != c.max || rl.Window != c.win {
			t.Errorf("%q: got %d/%s, want %d/%s", c.input, rl.Max, rl.Window, c.max, c.win)
		}
	}
}

func TestParseRateFlag_Invalid(t *testing.T) {
	cases := []string{
		"",
		"10",
		"abc/1s",
		"10/abc",
		"-1/1s",
		"10/0s",
		"10/-1s",
	}
	for _, c := range cases {
		_, err := ParseRateFlag(c)
		if err == nil {
			t.Errorf("%q: expected error, got nil", c)
		}
	}
}

func TestRateLimitFlag_String(t *testing.T) {
	rl := RateLimit{Max: 5, Window: time.Second}
	f := &rateLimitFlag{value: &rl}
	if got := f.String(); got != "5/1s" {
		t.Errorf("got %q, want \"5/1s\"", got)
	}
}

func TestRateLimitFlag_Type(t *testing.T) {
	f := &rateLimitFlag{value: &RateLimit{}}
	if f.Type() != "rate" {
		t.Errorf("expected type 'rate'")
	}
}
