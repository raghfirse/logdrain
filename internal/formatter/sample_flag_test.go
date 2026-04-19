package formatter

import "testing"

func TestParseSampleFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  uint64
	}{
		{"", 0},
		{"0", 0},
		{"1", 1},
		{"10", 10},
	}
	for _, tc := range cases {
		got, err := ParseSampleFlag(tc.input)
		if err != nil {
			t.Errorf("input %q: unexpected error %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %d want %d", tc.input, got, tc.want)
		}
	}
}

func TestParseSampleFlag_Invalid(t *testing.T) {
	cases := []string{"-1", "abc", "1.5"}
	for _, s := range cases {
		if _, err := ParseSampleFlag(s); err == nil {
			t.Errorf("expected error for input %q", s)
		}
	}
}

func TestSampleFlag_String_Default(t *testing.T) {
	f := &SampleFlag{}
	if f.String() != "0" {
		t.Errorf("expected '0', got %q", f.String())
	}
}

func TestSampleFlag_Set_Valid(t *testing.T) {
	f := &SampleFlag{}
	if err := f.Set("5"); err != nil {
		t.Fatal(err)
	}
	if f.Value != 5 {
		t.Errorf("expected 5, got %d", f.Value)
	}
}

func TestSampleFlag_Type(t *testing.T) {
	f := &SampleFlag{}
	if f.Type() != "uint" {
		t.Errorf("unexpected type %q", f.Type())
	}
}
