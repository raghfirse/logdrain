package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewTop_EmptyField(t *testing.T) {
	_, err := NewTop("", 5)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTop_ZeroN(t *testing.T) {
	_, err := NewTop("level", 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestTop_Apply_SuppressesLine(t *testing.T) {
	top, _ := NewTop("level", 3)
	out, emit := top.Apply(`{"level":"info"}`)
	if emit {
		t.Errorf("expected emit=false, got true")
	}
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestTop_Apply_NonJSONPassthrough(t *testing.T) {
	top, _ := NewTop("level", 3)
	out, emit := top.Apply("not json")
	if !emit {
		t.Errorf("expected emit=true for non-JSON")
	}
	if out != "not json" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTop_Flush_CountsCorrectly(t *testing.T) {
	top, _ := NewTop("level", 3)
	for i := 0; i < 5; i++ {
		top.Apply(`{"level":"error"}`)
	}
	for i := 0; i < 3; i++ {
		top.Apply(`{"level":"warn"}`)
	}
	top.Apply(`{"level":"info"}`)
	lines := top.Flush()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	var first map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if first["level"] != "error" {
		t.Errorf("expected top entry to be 'error', got %v", first["level"])
	}
	if int(first["count"].(float64)) != 5 {
		t.Errorf("expected count 5, got %v", first["count"])
	}
}

func TestTop_Flush_LimitedToN(t *testing.T) {
	top, _ := NewTop("status", 2)
	for _, v := range []string{"200", "404", "500", "200", "200", "404"} {
		top.Apply(`{"status":"` + v + `"}`)
	}
	lines := top.Flush()
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestTop_Reset_ClearsCounts(t *testing.T) {
	top, _ := NewTop("level", 3)
	top.Apply(`{"level":"info"}`)
	top.Reset()
	lines := top.Flush()
	if len(lines) != 0 {
		t.Errorf("expected 0 lines after reset, got %d", len(lines))
	}
}

func TestParseTopFlag_Valid(t *testing.T) {
	field, n, err := ParseTopFlag("level:10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if field != "level" || n != 10 {
		t.Errorf("got field=%q n=%d", field, n)
	}
}

func TestParseTopFlag_Invalid(t *testing.T) {
	if _, _, err := ParseTopFlag("level"); err == nil {
		t.Error("expected error for missing colon")
	}
	if _, _, err := ParseTopFlag("level:abc"); err == nil {
		t.Error("expected error for non-numeric n")
	}
}
