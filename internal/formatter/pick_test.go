package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewPick_Empty(t *testing.T) {
	p, err := NewPick(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, ok := p.Apply(`{"a":1,"b":2}`)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if out != `{"a":1,"b":2}` {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestPick_RetainsOnlySpecifiedKeys(t *testing.T) {
	p, _ := NewPick([]string{"msg", "level"})
	out, ok := p.Apply(`{"msg":"hello","level":"info","ts":"2024-01-01"}`)
	if !ok {
		t.Fatal("expected ok=true")
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected msg key")
	}
	if _, ok := m["level"]; !ok {
		t.Error("expected level key")
	}
	if _, ok := m["ts"]; ok {
		t.Error("expected ts key to be absent")
	}
}

func TestPick_CaseInsensitiveKey(t *testing.T) {
	p, _ := NewPick([]string{"MSG"})
	out, _ := p.Apply(`{"msg":"hello","level":"info"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected msg key (case-insensitive match)")
	}
	if _, ok := m["level"]; ok {
		t.Error("expected level to be absent")
	}
}

func TestPick_NonJSONPassthrough(t *testing.T) {
	p, _ := NewPick([]string{"msg"})
	raw := "not json at all"
	out, ok := p.Apply(raw)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if out != raw {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestPick_MissingKeyProducesEmptyObject(t *testing.T) {
	p, _ := NewPick([]string{"nonexistent"})
	out, _ := p.Apply(`{"msg":"hello"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty object, got %v", m)
	}
}

func TestParsePickFlag_CommaSeparated(t *testing.T) {
	keys, err := ParsePickFlag("msg, level, ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

func TestParsePickFlag_Empty(t *testing.T) {
	keys, err := ParsePickFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}
