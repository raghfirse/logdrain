package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewScript_Valid(t *testing.T) {
	_, err := NewScript([]string{`env=production`, `app={{.service}}`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewScript_InvalidRule(t *testing.T) {
	_, err := NewScript([]string{`noequalssign`})
	if err == nil {
		t.Fatal("expected error for rule without '='")
	}
}

func TestNewScript_EmptyKey(t *testing.T) {
	_, err := NewScript([]string{`=value`})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewScript_InvalidTemplate(t *testing.T) {
	_, err := NewScript([]string{`key={{.unclosed`})
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestScript_Apply_SetsLiteralField(t *testing.T) {
	s, _ := NewScript([]string{`env=production`})
	out := s.Apply(`{"msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["env"] != "production" {
		t.Errorf("expected env=production, got %v", obj["env"])
	}
}

func TestScript_Apply_TemplateField(t *testing.T) {
	s, _ := NewScript([]string{`summary={{.level}}: {{.msg}}`})
	out := s.Apply(`{"level":"error","msg":"boom"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["summary"] != "error: boom" {
		t.Errorf("unexpected summary: %v", obj["summary"])
	}
}

func TestScript_Apply_NonJSON(t *testing.T) {
	s, _ := NewScript([]string{`env=prod`})
	line := "not json at all"
	if got := s.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestScript_Apply_NoRules(t *testing.T) {
	s, _ := NewScript(nil)
	line := `{"msg":"hi"}`
	if got := s.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}
