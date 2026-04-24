package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewCounter_EmptyField(t *testing.T) {
	_, err := NewCounter("")
	if err == nil {
		t.Fatal("expected error for empty field, got nil")
	}
}

func TestNewCounter_ValidField(t *testing.T) {
	c, err := NewCounter("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestCounter_Apply_NonJSONSuppressed(t *testing.T) {
	c, _ := NewCounter("level")
	out := c.Apply("not json")
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
}

func TestCounter_Apply_SuppressesLine(t *testing.T) {
	c, _ := NewCounter("level")
	out := c.Apply(`{"level":"info","msg":"hello"}`)
	if out != "" {
		t.Fatalf("Apply should return empty string, got %q", out)
	}
}

func TestCounter_Flush_Empty(t *testing.T) {
	c, _ := NewCounter("level")
	lines := c.Flush()
	if len(lines) != 0 {
		t.Fatalf("expected no lines, got %d", len(lines))
	}
}

func TestCounter_Flush_CountsValues(t *testing.T) {
	c, _ := NewCounter("level")
	inputs := []string{
		`{"level":"info"}`,
		`{"level":"info"}`,
		`{"level":"error"}`,
	}
	for _, l := range inputs {
		c.Apply(l)
	}
	lines := c.Flush()
	if len(lines) != 2 {
		t.Fatalf("expected 2 summary lines, got %d", len(lines))
	}
	counts := map[string]int64{}
	for _, l := range lines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("invalid JSON in flush output: %v", err)
		}
		val := obj["level"].(string)
		counts[val] = int64(obj["count"].(float64))
	}
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", counts["error"])
	}
}

func TestCounter_Flush_ResetsState(t *testing.T) {
	c, _ := NewCounter("level")
	c.Apply(`{"level":"info"}`)
	c.Flush()
	lines := c.Flush()
	if len(lines) != 0 {
		t.Fatalf("expected empty flush after reset, got %d lines", len(lines))
	}
}

func TestParseCountFlag_Valid(t *testing.T) {
	field, err := ParseCountFlag("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if field != "level" {
		t.Errorf("expected 'level', got %q", field)
	}
}

func TestParseCountFlag_Empty(t *testing.T) {
	_, err := ParseCountFlag("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}
