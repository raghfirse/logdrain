package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewRenamer_NoRules(t *testing.T) {
	r := NewRenamer(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %s", got)
	}
}

func TestRenamer_RenamesField(t *testing.T) {
	r := NewRenamer([]RenameRule{{From: "msg", To: "message"}})
	got := r.Apply(`{"level":"info","msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["msg"]; ok {
		t.Error("old key 'msg' should be removed")
	}
	if obj["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", obj["message"])
	}
}

func TestRenamer_CaseInsensitiveFrom(t *testing.T) {
	r := NewRenamer([]RenameRule{{From: "MSG", To: "message"}})
	got := r.Apply(`{"msg":"hi"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["message"] != "hi" {
		t.Errorf("expected message=hi, got %v", obj["message"])
	}
}

func TestRenamer_MissingKeyNoOp(t *testing.T) {
	r := NewRenamer([]RenameRule{{From: "missing", To: "other"}})
	line := `{"level":"warn"}`
	got := r.Apply(line)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["other"]; ok {
		t.Error("'other' should not be present")
	}
}

func TestRenamer_NonJSONPassthrough(t *testing.T) {
	r := NewRenamer([]RenameRule{{From: "a", To: "b"}})
	line := "plain text line"
	if got := r.Apply(line); got != line {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestParseRenameFlag_Valid(t *testing.T) {
	rules, err := ParseRenameFlag([]string{"msg:message", "ts:timestamp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "msg" || rules[0].To != "message" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRenameFlag_Invalid(t *testing.T) {
	cases := []string{"nocolon", ":nokey", "novalue:"}
	for _, c := range cases {
		_, err := ParseRenameFlag([]string{c})
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestRenameFlagValue_SetAndRules(t *testing.T) {
	var f RenameFlagValue
	if err := f.Set("msg:message"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rules, err := f.Rules()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].From != "msg" {
		t.Errorf("unexpected rules: %+v", rules)
	}
}

func TestRenameFlagValue_SetInvalid(t *testing.T) {
	var f RenameFlagValue
	if err := f.Set("badvalue"); err == nil {
		t.Error("expected error for invalid value")
	}
}
