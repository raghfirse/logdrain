package formatter

import (
	"testing"
)

func TestNewOmit_Empty(t *testing.T) {
	_, err := NewOmit([]string{})
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestOmit_RemovesField(t *testing.T) {
	o, err := NewOmit([]string{"secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"msg":"hello","secret":"abc"}`
	out, keep := o.Apply(line)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if contains(out, "secret") {
		t.Errorf("expected 'secret' to be removed, got: %s", out)
	}
	if !contains(out, "msg") {
		t.Errorf("expected 'msg' to remain, got: %s", out)
	}
}

func TestOmit_CaseInsensitiveKey(t *testing.T) {
	o, _ := NewOmit([]string{"Secret"})
	line := `{"secret":"value","keep":"yes"}`
	out, _ := o.Apply(line)
	if contains(out, "secret") {
		t.Errorf("expected 'secret' to be removed, got: %s", out)
	}
}

func TestOmit_NonJSONPassthrough(t *testing.T) {
	o, _ := NewOmit([]string{"x"})
	line := "plain text log"
	out, keep := o.Apply(line)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if out != line {
		t.Errorf("expected unchanged line, got: %s", out)
	}
}

func TestOmit_MultipleKeys(t *testing.T) {
	o, _ := NewOmit([]string{"a", "b"})
	line := `{"a":1,"b":2,"c":3}`
	out, _ := o.Apply(line)
	if contains(out, `"a"`) || contains(out, `"b"`) {
		t.Errorf("expected a and b removed, got: %s", out)
	}
	if !contains(out, `"c"`) {
		t.Errorf("expected c to remain, got: %s", out)
	}
}

func TestOmit_MissingKeyNoOp(t *testing.T) {
	o, _ := NewOmit([]string{"missing"})
	line := `{"present":"yes"}`
	out, _ := o.Apply(line)
	if !contains(out, "present") {
		t.Errorf("expected 'present' to remain, got: %s", out)
	}
}

func TestParseOmitFlag_Valid(t *testing.T) {
	keys, err := ParseOmitFlag("foo,bar,baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

func TestParseOmitFlag_Empty(t *testing.T) {
	_, err := ParseOmitFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
