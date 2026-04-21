package formatter

import (
	"testing"
)

func TestDiff_NonJSONPassthrough(t *testing.T) {
	d := NewDiff(nil)
	line := "not json"
	if got := d.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestDiff_FirstLineReturnedAsIs(t *testing.T) {
	d := NewDiff(nil)
	line := `{"level":"info","msg":"start"}`
	if got := d.Apply(line); got != line {
		t.Errorf("expected first line unchanged, got %q", got)
	}
}

func TestDiff_UnchangedFieldsOmitted(t *testing.T) {
	d := NewDiff(nil)
	d.Apply(`{"level":"info","status":200}`)
	out := d.Apply(`{"level":"info","status":200}`)
	if out == `{"level":"info","status":200}` {
		t.Error("expected diff output, not original line")
	}
	if !contains(out, `"_diff"`) {
		t.Errorf("expected _diff key in output, got %q", out)
	}
	if contains(out, `"level"`) {
		t.Errorf("unchanged field 'level' should be omitted, got %q", out)
	}
}

func TestDiff_ChangedFieldIncluded(t *testing.T) {
	d := NewDiff(nil)
	d.Apply":200}`)
	out := d.Apply(`{"status":500}`)
	if !contains(out, `"status"`) {
		t.Errorf("expected changed field 'status' in output, got %q", out)
	}
	if !contains(out, `500`) {
		t.Errorf("expected new value 500 in output, got %q", out)
	}
}

func TestDiff_RemovedFieldIsNull(t *testing.T) {
	d := NewDiff(nil)
	d.Apply(`{"status":200,"extra":"yes"}`)
	out := d.Apply(`{"status":200}`)
	if !contains(out, `"extra":null`) {
		t.Errorf("expected removed field to be null, got %q", out)
	}
}

func TestDiff_SpecificFields(t *testing.T) {
	d := NewDiff([]string{"status"})
	d.Apply(`{"status":200,"msg":"a"}`)
	out := d.Apply(`{"status":200,"msg":"b"}`)
	// msg changed but is not in the watched fields
	if contains(out, `"msg"`) {
		t.Errorf("non-watched field 'msg' should be absent, got %q", out)
	}
}

func TestDiff_Reset(t *testing.T) {
	d := NewDiff(nil)
	line := `{"x":1}`
	d.Apply(line)
	d.Reset()
	// After reset the next line should be treated as first
	if got := d.Apply(line); got != line {
		t.Errorf("expected first-line passthrough after reset, got %q", got)
	}
}

func TestParseDiffFlag_Empty(t *testing.T) {
	fields, err := ParseDiffFlag("")
	if err != nil || fields != nil {
		t.Errorf("expected nil fields and no error, got %v, %v", fields, err)
	}
}

func TestParseDiffFlag_Multiple(t *testing.T) {
	fields, err := ParseDiffFlag("a,b,c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 3 || fields[0] != "a" || fields[2] != "c" {
		t.Errorf("unexpected fields: %v", fields)
	}
}

func TestParseDiffFlag_Invalid(t *testing.T) {
	_, err := ParseDiffFlag("a,,b")
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestDiffFlagValue_SetAndGet(t *testing.T) {
	f := NewDiffFlagValue()
	if err := f.Set("x,y"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Fields()) != 2 {
		t.Errorf("expected 2 fields, got %v", f.Fields())
	}
	if f.String() != "x,y" {
		t.Errorf("unexpected string: %q", f.String())
	}
}

func TestDiffFlagValue_Type(t *testing.T) {
	f := NewDiffFlagValue()
	if f.Type() != "fields" {
		t.Errorf("unexpected type: %q", f.Type())
	}
}

// contains is a helper used across formatter tests.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
