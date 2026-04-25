package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewRollup_EmptyField(t *testing.T) {
	_, err := NewRollup("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRollup_ValidField(t *testing.T) {
	r, err := NewRollup("service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rollup")
	}
}

func TestRollup_Apply_NonJSONPassthrough(t *testing.T) {
	r, _ := NewRollup("service")
	out := r.Apply("not json")
	if out != "not json" {
		t.Errorf("expected passthrough, got %q", out)
	}
}

func TestRollup_Apply_SuppressesJSONLine(t *testing.T) {
	r, _ := NewRollup("service")
	out := r.Apply(`{"service":"api","msg":"hello"}`)
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestRollup_Flush_CountsCorrectly(t *testing.T) {
	r, _ := NewRollup("service")
	r.Apply(`{"service":"api"}`)
	r.Apply(`{"service":"api"}`)
	r.Apply(`{"service":"worker"}`)

	lines := r.Flush()
	if len(lines) != 2 {
		t.Fatalf("expected 2 summary lines, got %d", len(lines))
	}

	counts := map[string]int{}
	for _, l := range lines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("invalid JSON in flush output: %v", err)
		}
		svc := obj["service"].(string)
		counts[svc] = int(obj["count"].(float64))
	}
	if counts["api"] != 2 {
		t.Errorf("expected api count=2, got %d", counts["api"])
	}
	if counts["worker"] != 1 {
		t.Errorf("expected worker count=1, got %d", counts["worker"])
	}
}

func TestRollup_Flush_ResetsState(t *testing.T) {
	r, _ := NewRollup("service")
	r.Apply(`{"service":"api"}`)
	r.Flush()
	lines := r.Flush()
	if len(lines) != 0 {
		t.Errorf("expected empty flush after reset, got %d lines", len(lines))
	}
}

func TestRollup_MissingField(t *testing.T) {
	r, _ := NewRollup("service")
	r.Apply(`{"msg":"no service field"}`)
	lines := r.Flush()
	if len(lines) != 1 {
		t.Fatalf("expected 1 summary line, got %d", len(lines))
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(lines[0]), &obj)
	if obj["service"] != "(missing)" {
		t.Errorf("expected (missing) key, got %v", obj["service"])
	}
}

func TestParseRollupFlag_Valid(t *testing.T) {
	v, err := ParseRollupFlag("service")
	if err != nil || v != "service" {
		t.Errorf("unexpected result: %q %v", v, err)
	}
}

func TestParseRollupFlag_Empty(t *testing.T) {
	_, err := ParseRollupFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}
