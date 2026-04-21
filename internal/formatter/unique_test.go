package formatter

import (
	"testing"
)

func TestNewUnique_EmptyField(t *testing.T) {
	_, err := NewUnique("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewUnique_ValidField(t *testing.T) {
	u, err := NewUnique("request_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil {
		t.Fatal("expected non-nil Unique")
	}
}

func TestUnique_FirstOccurrenceNotDuplicate(t *testing.T) {
	u, _ := NewUnique("id")
	line := `{"id":"abc","msg":"hello"}`
	if u.IsDuplicate(line) {
		t.Error("first occurrence should not be a duplicate")
	}
}

func TestUnique_SecondOccurrenceIsDuplicate(t *testing.T) {
	u, _ := NewUnique("id")
	line := `{"id":"abc","msg":"hello"}`
	u.IsDuplicate(line)
	if !u.IsDuplicate(line) {
		t.Error("second occurrence should be a duplicate")
	}
}

func TestUnique_DifferentValuesNotDuplicate(t *testing.T) {
	u, _ := NewUnique("id")
	u.IsDuplicate(`{"id":"abc"}`)
	if u.IsDuplicate(`{"id":"xyz"}`) {
		t.Error("different field value should not be a duplicate")
	}
}

func TestUnique_NonJSONPassthrough(t *testing.T) {
	u, _ := NewUnique("id")
	if u.IsDuplicate("not json at all") {
		t.Error("non-JSON line should never be a duplicate")
	}
	if u.IsDuplicate("not json at all") {
		t.Error("repeated non-JSON line should still not be a duplicate")
	}
}

func TestUnique_MissingFieldNotDuplicate(t *testing.T) {
	u, _ := NewUnique("id")
	line := `{"msg":"no id field here"}`
	u.IsDuplicate(line)
	if u.IsDuplicate(line) {
		t.Error("lines missing the target field should never be duplicates")
	}
}

func TestUnique_CaseInsensitiveKey(t *testing.T) {
	u, _ := NewUnique("RequestID")
	if u.IsDuplicate(`{"requestid":"abc"}`) {
		t.Error("first occurrence should not be duplicate")
	}
	if !u.IsDuplicate(`{"REQUESTID":"abc"}`) {
		t.Error("same value under different case key should be duplicate")
	}
}

func TestUnique_Reset(t *testing.T) {
	u, _ := NewUnique("id")
	line := `{"id":"abc"}`
	u.IsDuplicate(line)
	u.Reset()
	if u.IsDuplicate(line) {
		t.Error("after reset, line should not be a duplicate")
	}
}

func TestParseUniqueFlag_Valid(t *testing.T) {
	field, err := ParseUniqueFlag("  request_id  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if field != "request_id" {
		t.Errorf("expected 'request_id', got %q", field)
	}
}

func TestParseUniqueFlag_Empty(t *testing.T) {
	_, err := ParseUniqueFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}
