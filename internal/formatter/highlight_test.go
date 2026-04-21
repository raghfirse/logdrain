package formatter

import (
	"strings"
	"testing"
)

func TestHighlightFields_NoHighlights(t *testing.T) {
	fields := map[string]interface{}{"msg": "hello", "level": "info"}
	result := HighlightFields(fields, nil)
	if result["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", result["msg"])
	}
}

func TestHighlightFields_AppliesColor(t *testing.T) {
	fields := map[string]interface{}{"level": "error", "msg": "oops"}
	highlights := []FieldHighlight{{Field: "level", Color: "\033[31m"}}
	result := HighlightFields(fields, highlights)
	val, ok := result["level"].(string)
	if !ok {
		t.Fatal("expected string value")
	}
	if !strings.Contains(val, "\033[31m") {
		t.Errorf("expected ANSI red in value, got %q", val)
	}
	if !strings.Contains(val, "error") {
		t.Errorf("expected original value preserved, got %q", val)
	}
}

func TestHighlightFields_CaseInsensitiveKey(t *testing.T) {
	fields := map[string]interface{}{"Level": "warn"}
	highlights := []FieldHighlight{{Field: "level", Color: "\033[33m"}}
	result := HighlightFields(fields, highlights)
	val, _ := result["Level"].(string)
	if !strings.Contains(val, "\033[33m") {
		t.Errorf("expected color applied case-insensitively, got %q", val)
	}
}

func TestHighlightFields_NonStringValue(t *testing.T) {
	// Non-string field values should be left unmodified.
	fields := map[string]interface{}{"count": 42}
	highlights := []FieldHighlight{{Field: "count", Color: "\033[31m"}}
	result := HighlightFields(fields, highlights)
	if result["count"] != 42 {
		t.Errorf("expected non-string value to be unchanged, got %v", result["count"])
	}
}

func TestHighlightFields_ResetCodeAppended(t *testing.T) {
	// Ensure the ANSI reset code is appended after the colored value so
	// subsequent output is not unintentionally colored.
	fields := map[string]interface{}{"level": "info"}
	highlights := []FieldHighlight{{Field: "level", Color: "\033[34m"}}
	result := HighlightFields(fields, highlights)
	val, ok := result["level"].(string)
	if !ok {
		t.Fatal("expected string value")
	}
	if !strings.Contains(val, "\033[0m") {
		t.Errorf("expected ANSI reset code in value, got %q", val)
	}
}

func TestParseHighlightFlag_Valid(t *testing.T) {
	hl := ParseHighlightFlag([]string{"level=red", "msg=cyan"})
	if len(hl) != 2 {
		t.Fatalf("expected 2 highlights, got %d", len(hl))
	}
	if hl[0].Field != "level" {
		t.Errorf("expected field 'level', got %q", hl[0].Field)
	}
	if hl[0].Color != "\033[31m" {
		t.Errorf("expected red ANSI, got %q", hl[0].Color)
	}
}

func TestParseHighlightFlag_InvalidExpr(t *testing.T) {
	hl := ParseHighlightFlag([]string{"nodequals"})
	if len(hl) != 0 {
		t.Errorf("expected 0 highlights for invalid expr, got %d", len(hl))
	}
}

func TestParseHighlightFlag_UnknownColor_DefaultsYellow(t *testing.T) {
	hl := ParseHighlightFlag([]string{"field=neonpink"})
	if len(hl) != 1 {
		t.Fatalf("expected 1 highlight")
	}
	if hl[0].Color != "\033[33m" {
		t.Errorf("expected default yellow, got %q", hl[0].Color)
	}
}
