package formatter

import (
	"testing"
)

func TestNewContextFields_Extract_Match(t *testing.T) {
	cf := NewContextFields([]string{"request_id", "user"})
	line := []byte(`{"level":"info","request_id":"abc123","user":"alice","msg":"ok"}`)
	got := cf.Extract(line)
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got["request_id"] != "abc123" {
		t.Errorf("expected request_id=abc123, got %v", got["request_id"])
	}
	if got["user"] != "alice" {
		t.Errorf("expected user=alice, got %v", got["user"])
	}
	if _, ok := got["level"]; ok {
		t.Error("level should not be extracted")
	}
}

func TestNewContextFields_Extract_NoMatch(t *testing.T) {
	cf := NewContextFields([]string{"trace_id"})
	line := []byte(`{"level":"info","msg":"hello"}`)
	if got := cf.Extract(line); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNewContextFields_Extract_InvalidJSON(t *testing.T) {
	cf := NewContextFields([]string{"x"})
	if got := cf.Extract([]byte("not json")); got != nil {
		t.Errorf("expected nil for invalid JSON, got %v", got)
	}
}

func TestNewContextFields_Extract_EmptyKeys(t *testing.T) {
	cf := NewContextFields(nil)
	line := []byte(`{"a":"1"}`)
	if got := cf.Extract(line); got != nil {
		t.Errorf("expected nil with no keys configured, got %v", got)
	}
}

func TestNewContextFields_CaseInsensitiveKey(t *testing.T) {
	cf := NewContextFields([]string{"RequestID"})
	line := []byte(`{"requestid":"xyz"}`)
	got := cf.Extract(line)
	if got == nil || got["requestid"] != "xyz" {
		t.Errorf("expected case-insensitive match, got %v", got)
	}
}

func TestParseContextFlag_Empty(t *testing.T) {
	if got := ParseContextFlag(""); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestParseContextFlag_Multiple(t *testing.T) {
	got := ParseContextFlag("request_id, user , trace")
	if len(got) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(got))
	}
	if got[1] != "user" {
		t.Errorf("expected trimmed 'user', got %q", got[1])
	}
}
