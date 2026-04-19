package filter

import (
	"testing"
)

func TestNew_ValidExpressions(t *testing.T) {
	f, err := New([]string{"level=info", "service=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(f.rules))
	}
}

func TestNew_InvalidExpression(t *testing.T) {
	_, err := New([]string{"badexpr"})
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestMatch_NoRules(t *testing.T) {
	f, _ := New(nil)
	if !f.Match([]byte(`{"level":"error"}`)) {
		t.Fatal("expected match with no rules")
	}
}

func TestMatch_SingleRule_Match(t *testing.T) {
	f, _ := New([]string{"level=info"})
	if !f.Match([]byte(`{"level":"info","msg":"started"}`)) {
		t.Fatal("expected match")
	}
}

func TestMatch_SingleRule_NoMatch(t *testing.T) {
	f, _ := New([]string{"level=info"})
	if f.Match([]byte(`{"level":"error","msg":"failed"}`)) {
		t.Fatal("expected no match")
	}
}

func TestMatch_MultipleRules_AllMatch(t *testing.T) {
	f, _ := New([]string{"level=info", "service=api"})
	if !f.Match([]byte(`{"level":"info","service":"api","msg":"ok"}`)) {
		t.Fatal("expected match")
	}
}

func TestMatch_MultipleRules_PartialMatch(t *testing.T) {
	f, _ := New([]string{"level=info", "service=api"})
	if f.Match([]byte(`{"level":"info","service":"worker"}`)) {
		t.Fatal("expected no match when only one rule matches")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f, _ := New([]string{"level=info"})
	if f.Match([]byte(`not json`)) {
		t.Fatal("expected no match for invalid JSON")
	}
}

func TestMatch_MissingField(t *testing.T) {
	f, _ := New([]string{"level=info"})
	if f.Match([]byte(`{"msg":"hello"}`)) {
		t.Fatal("expected no match when field is absent")
	}
}
