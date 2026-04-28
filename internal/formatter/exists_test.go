package formatter

import (
	"testing"
)

func TestNewExists_NoFields(t *testing.T) {
	e, err := NewExists(ExistsConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Exists")
	}
}

func TestExists_Apply_NonJSONPassthrough(t *testing.T) {
	e, _ := NewExists(ExistsConfig{Required: []string{"user"}})
	line := "not json at all"
	if got := e.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestExists_Apply_NoRulesPassthrough(t *testing.T) {
	e, _ := NewExists(ExistsConfig{})
	line := `{"msg":"hello"}`
	if got := e.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestExists_Apply_RequiredPresent(t *testing.T) {
	e, _ := NewExists(ExistsConfig{Required: []string{"user"}})
	line := `{"user":"alice","msg":"hello"}`
	if got := e.Apply(line); got != line {
		t.Errorf("expected line returned, got %q", got)
	}
}

func TestExists_Apply_RequiredMissing(t *testing.T) {
	e, _ := NewExists(ExistsConfig{Required: []string{"user"}})
	line := `{"msg":"hello"}`
	if got := e.Apply(line); got != "" {
		t.Errorf("expected suppression, got %q", got)
	}
}

func TestExists_Apply_ForbiddenAbsent(t *testing.T) {
	e, _ := NewExists(ExistsConfig{Forbidden: []string{"debug"}})
	line := `{"msg":"hello"}`
	if got := e.Apply(line); got != line {
		t.Errorf("expected line returned, got %q", got)
	}
}

func TestExists_Apply_ForbiddenPresent(t *testing.T) {
	e, _ := NewExists(ExistsConfig{Forbidden: []string{"debug"}})
	line := `{"debug":true,"msg":"hello"}`
	if got := e.Apply(line); got != "" {
		t.Errorf("expected suppression, got %q", got)
	}
}

func TestExists_Apply_CaseInsensitiveKey(t *testing.T) {
	e, _ := NewExists(ExistsConfig{Required: []string{"User"}})
	line := `{"user":"alice"}`
	if got := e.Apply(line); got != line {
		t.Errorf("expected case-insensitive match, got %q", got)
	}
}

func TestParseExistsFlag_Empty(t *testing.T) {
	cfg, err := ParseExistsFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Required) != 0 || len(cfg.Forbidden) != 0 {
		t.Errorf("expected empty config, got %+v", cfg)
	}
}

func TestParseExistsFlag_RequiredOnly(t *testing.T) {
	cfg, err := ParseExistsFlag("user,host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Required) != 2 || cfg.Required[0] != "user" || cfg.Required[1] != "host" {
		t.Errorf("unexpected required: %v", cfg.Required)
	}
	if len(cfg.Forbidden) != 0 {
		t.Errorf("unexpected forbidden: %v", cfg.Forbidden)
	}
}

func TestParseExistsFlag_ForbiddenOnly(t *testing.T) {
	cfg, err := ParseExistsFlag("!debug,!trace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Forbidden) != 2 || cfg.Forbidden[0] != "debug" || cfg.Forbidden[1] != "trace" {
		t.Errorf("unexpected forbidden: %v", cfg.Forbidden)
	}
}

func TestParseExistsFlag_Mixed(t *testing.T) {
	cfg, err := ParseExistsFlag("user,!debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Required) != 1 || cfg.Required[0] != "user" {
		t.Errorf("unexpected required: %v", cfg.Required)
	}
	if len(cfg.Forbidden) != 1 || cfg.Forbidden[0] != "debug" {
		t.Errorf("unexpected forbidden: %v", cfg.Forbidden)
	}
}
