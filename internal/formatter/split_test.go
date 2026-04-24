package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewSplit_EmptyField(t *testing.T) {
	_, err := NewSplit("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSplit_ValidField(t *testing.T) {
	s, err := NewSplit("items")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Split")
	}
}

func TestSplit_NonJSONPassthrough(t *testing.T) {
	s, _ := NewSplit("items")
	result := s.Apply("not json")
	if len(result) != 1 || result[0] != "not json" {
		t.Fatalf("expected passthrough, got %v", result)
	}
}

func TestSplit_MissingFieldPassthrough(t *testing.T) {
	s, _ := NewSplit("items")
	result := s.Apply(`{"level":"info"}`)
	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}
}

func TestSplit_NonArrayFieldPassthrough(t *testing.T) {
	s, _ := NewSplit("items")
	result := s.Apply(`{"items":"not-an-array"}`)
	if len(result) != 1 {
		t.Fatalf("expected 1 line passthrough, got %d", len(result))
	}
}

func TestSplit_ArrayField_EmitsOnePerElement(t *testing.T) {
	s, _ := NewSplit("items")
	result := s.Apply(`{"level":"info","items":["a","b","c"]}`)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	for i, want := range []string{`"a"`, `"b"`, `"c"`} {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(result[i]), &obj); err != nil {
			t.Fatalf("line %d not valid JSON: %v", i, err)
		}
		if string(obj["items"]) != want {
			t.Errorf("line %d: items = %s, want %s", i, obj["items"], want)
		}
		if string(obj["level"]) != `"info"` {
			t.Errorf("line %d: level missing or wrong", i)
		}
	}
}

func TestSplit_CaseInsensitiveField(t *testing.T) {
	s, _ := NewSplit("Items")
	result := s.Apply(`{"items":[1,2]}`)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}
}

func TestSplit_EmptyArray(t *testing.T) {
	s, _ := NewSplit("items")
	result := s.Apply(`{"items":[]}`)
	if len(result) != 0 {
		t.Fatalf("expected 0 lines for empty array, got %d", len(result))
	}
}

func TestParseSplitFlag_Valid(t *testing.T) {
	s, err := ParseSplitFlag("tags")
	if err != nil || s == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseSplitFlag_Invalid(t *testing.T) {
	_, err := ParseSplitFlag("  ")
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}
