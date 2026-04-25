package formatter

import (
	"testing"
)

func TestNewMask_NoSpecs(t *testing.T) {
	m, err := NewMask(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Apply(`{"card":"1234567890121234"}`); got != `{"card":"1234567890121234"}` {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestNewMask_InvalidSpec(t *testing.T) {
	_, err := NewMask([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestNewMask_InvalidPattern(t *testing.T) {
	_, err := NewMask([]string{"field:[invalid:mask"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMask_Apply_NonJSON(t *testing.T) {
	m, _ := NewMask([]string{`card:\d{12}(\d{4}):****$1`})
	line := "not json at all"
	if got := m.Apply(line); got != line {
		t.Errorf("expected passthrough for non-JSON, got %s", got)
	}
}

func TestMask_Apply_MasksMatchingField(t *testing.T) {
	m, err := NewMask([]string{`card:\d{12}(\d{4}):****$1`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := `{"card":"1234567890121234"}`
	out := m.Apply(input)
	if out == input {
		t.Errorf("expected field to be masked, got original: %s", out)
	}
	if !contains(out, "****1234") {
		t.Errorf("expected masked value ****1234 in output, got: %s", out)
	}
}

func TestMask_Apply_NoMatchingField(t *testing.T) {
	m, _ := NewMask([]string{`card:\d{12}(\d{4}):****$1`})
	input := `{"email":"user@example.com"}`
	if got := m.Apply(input); got != input {
		t.Errorf("expected passthrough for non-matching field, got %s", got)
	}
}

func TestMask_Apply_CaseInsensitiveKey(t *testing.T) {
	m, err := NewMask([]string{`card:\d{12}(\d{4}):****$1`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := `{"Card":"1234567890121234"}`
	out := m.Apply(input)
	if !contains(out, "****1234") {
		t.Errorf("expected case-insensitive match, got: %s", out)
	}
}

func TestMask_Apply_PatternNoMatch(t *testing.T) {
	m, _ := NewMask([]string{`card:\d{16}:REDACTED`})
	input := `{"card":"short"}`
	if got := m.Apply(input); got != input {
		t.Errorf("expected passthrough when pattern does not match, got %s", got)
	}
}

func TestMaskFlagValue_Set_Valid(t *testing.T) {
	var specs []string
	fv := NewMaskFlagValue(&specs)
	if err := fv.Set(`card:\d{12}(\d{4}):****$1`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 1 {
		t.Errorf("expected 1 spec, got %d", len(specs))
	}
}

func TestMaskFlagValue_Set_Invalid(t *testing.T) {
	var specs []string
	fv := NewMaskFlagValue(&specs)
	if err := fv.Set("noColons"); err == nil {
		t.Fatal("expected error for invalid spec")
	}
}

func TestMaskFlagValue_String_Default(t *testing.T) {
	var specs []string
	fv := NewMaskFlagValue(&specs)
	if got := fv.String(); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

// contains is a local helper (avoids import of strings in test file).
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
