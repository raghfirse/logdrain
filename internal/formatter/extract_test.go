package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewExtractor_Empty(t *testing.T) {
	e, err := NewExtractor(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(e.rules) != 0 {
		t.Errorf("expected no rules, got %d", len(e.rules))
	}
}

func TestNewExtractor_InvalidPath(t *testing.T) {
	_, err := NewExtractor([]string{"a..b"})
	if err == nil {
		t.Fatal("expected error for empty path segment")
	}
}

func TestNewExtractor_ValidPath(t *testing.T) {
	e, err := NewExtractor([]string{"meta.request_id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(e.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(e.rules))
	}
	if e.rules[0].alias != "request_id" {
		t.Errorf("expected alias 'request_id', got %q", e.rules[0].alias)
	}
}

func TestNewExtractor_Alias(t *testing.T) {
	e, err := NewExtractor([]string{"meta.request_id->rid"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.rules[0].alias != "rid" {
		t.Errorf("expected alias 'rid', got %q", e.rules[0].alias)
	}
}

func TestExtract_Apply_NonJSON(t *testing.T) {
	e, _ := NewExtractor([]string{"a.b"})
	out := e.Apply("not json")
	if out != "not json" {
		t.Errorf("expected passthrough, got %q", out)
	}
}

func TestExtract_Apply_NestedField(t *testing.T) {
	e, _ := NewExtractor([]string{"meta.request_id"})
	input := `{"level":"info","meta":{"request_id":"abc123","user":"bob"}}`
	out := e.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["request_id"] != "abc123" {
		t.Errorf("expected request_id=abc123, got %v", obj["request_id"])
	}
}

func TestExtract_Apply_MissingPath(t *testing.T) {
	e, _ := NewExtractor([]string{"meta.missing_key"})
	input := `{"meta":{"other":"val"}}`
	out := e.Apply(input)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if _, ok := obj["missing_key"]; ok {
		t.Error("expected missing_key to not be present")
	}
}

func TestExtract_Apply_WithAlias(t *testing.T) {
	e, _ := NewExtractor([]string{"meta.request_id->rid"})
	input := `{"meta":{"request_id":"xyz"}}`
	out := e.Apply(input)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["rid"] != "xyz" {
		t.Errorf("expected rid=xyz, got %v", obj["rid"])
	}
	if _, ok := obj["request_id"]; ok {
		t.Error("expected no 'request_id' top-level key")
	}
}

func TestExtract_Apply_CaseInsensitive(t *testing.T) {
	e, _ := NewExtractor([]string{"Meta.RequestID"})
	input := `{"meta":{"requestid":"hello"}}`
	out := e.Apply(input)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["RequestID"] != "hello" {
		t.Errorf("expected RequestID=hello, got %v", obj["RequestID"])
	}
}
