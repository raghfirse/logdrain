package formatter

import (
	"testing"
)

func TestNewTemplate_Valid(t *testing.T) {
	_, err := NewTemplate(`{{index . "level"}} {{index . "msg"}}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewTemplate_Empty(t *testing.T) {
	_, err := NewTemplate("")
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestNewTemplate_Invalid(t *testing.T) {
	_, err := NewTemplate("{{.Unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestTemplate_Render_ValidJSON(t *testing.T) {
	tmpl, err := NewTemplate(`[{{index . "level"}}] {{index . "msg"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"info","msg":"hello world"}`
	got := tmpl.Render(line)
	want := "[info] hello world"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestTemplate_Render_InvalidJSON(t *testing.T) {
	tmpl, _ := NewTemplate(`{{index . "msg"}}`)
	line := "not json at all"
	got := tmpl.Render(line)
	if got != line {
		t.Errorf("Render() = %q, want original %q", got, line)
	}
}

func TestTemplate_Render_MissingKey(t *testing.T) {
	tmpl, _ := NewTemplate(`{{index . "nonexistent"}}`)
	line := `{"level":"debug"}`
	got := tmpl.Render(line)
	// missingkey=zero: missing key renders as empty string / zero value
	if got == line {
		t.Errorf("expected rendered output, got original line back")
	}
}

func TestParseTemplateFlag_Valid(t *testing.T) {
	_, err := ParseTemplateFlag(`{{index . "level"}}: {{index . "msg"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseTemplateFlag_Invalid(t *testing.T) {
	_, err := ParseTemplateFlag("{{bad")
	if err == nil {
		t.Fatal("expected error for bad template")
	}
}
