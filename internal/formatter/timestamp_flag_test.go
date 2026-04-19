package formatter

import "testing"

func TestParseTimestampMode_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  TimestampMode
	}{
		{"off", TimestampOff},
		{"none", TimestampOff},
		{"", TimestampOff},
		{"short", TimestampShort},
		{"full", TimestampFull},
		{"long", TimestampFull},
	}
	for _, tc := range cases {
		got, err := ParseTimestampMode(tc.input)
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("input %q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseTimestampMode_Invalid(t *testing.T) {
	_, err := ParseTimestampMode("iso")
	if err == nil {
		t.Error("expected error for unknown mode")
	}
}

func TestTimestampMode_String(t *testing.T) {
	if TimestampOff.String() != "off" {
		t.Errorf("expected 'off'")
	}
	if TimestampShort.String() != "short" {
		t.Errorf("expected 'short'")
	}
	if TimestampFull.String() != "full" {
		t.Errorf("expected 'full'")
	}
}
