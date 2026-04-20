package formatter

import (
	"testing"
)

func TestParseAggregateFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		field string
		mode  AggregateMode
	}{
		{"latency:count", "latency", AggregateCount},
		{"bytes:sum", "bytes", AggregateSum},
		{"duration:min", "duration", AggregateMin},
		{"score:max", "score", AggregateMax},
		{"val:COUNT", "val", AggregateCount},
		{"val:Sum", "val", AggregateSum},
	}
	for _, tc := range cases {
		f, m, err := ParseAggregateFlag(tc.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", tc.input, err)
			continue
		}
		if f != tc.field {
			t.Errorf("%q: field = %q, want %q", tc.input, f, tc.field)
		}
		if m != tc.mode {
			t.Errorf("%q: mode = %v, want %v", tc.input, m, tc.mode)
		}
	}
}

func TestParseAggregateFlag_Invalid(t *testing.T) {
	cases := []string{
		"",
		"nocodon",
		":count",
		"field:unknown",
		"field:",
	}
	for _, tc := range cases {
		_, _, err := ParseAggregateFlag(tc)
		if err == nil {
			t.Errorf("%q: expected error, got nil", tc)
		}
	}
}

func TestAggregateFlagValue_Set_Valid(t *testing.T) {
	var f AggregateFlagValue
	if err := f.Set("latency:sum"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "latency" {
		t.Errorf("field = %q, want latency", f.Field)
	}
	if f.Mode != AggregateSum {
		t.Errorf("mode = %v, want AggregateSum", f.Mode)
	}
}

func TestAggregateFlagValue_String_Default(t *testing.T) {
	var f AggregateFlagValue
	if f.String() != "" {
		t.Errorf("expected empty string for unset flag")
	}
}

func TestAggregateFlagValue_String_AfterSet(t *testing.T) {
	var f AggregateFlagValue
	_ = f.Set("score:max")
	if f.String() != "score:max" {
		t.Errorf("String() = %q, want score:max", f.String())
	}
}

func TestAggregateFlagValue_Type(t *testing.T) {
	var f AggregateFlagValue
	if f.Type() != "field:mode" {
		t.Errorf("Type() = %q, want field:mode", f.Type())
	}
}
