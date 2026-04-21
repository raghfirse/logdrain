package formatter

import (
	"testing"
)

func TestNewGrep_ValidPattern(t *testing.T) {
	g, err := NewGrep(`error`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil Grep")
	}
}

func TestNewGrep_InvalidPattern(t *testing.T) {
	_, err := NewGrep(`[invalid`, nil)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestGrep_Match_WholeLine(t *testing.T) {
	g, _ := NewGrep(`error`, nil)
	if !g.Match(`{"level":"error","msg":"boom"}`) {
		t.Error("expected match")
	}
	if g.Match(`{"level":"info","msg":"ok"}`) {
		t.Error("expected no match")
	}
}

func TestGrep_Match_SpecificField(t *testing.T) {
	g, _ := NewGrep(`boom`, []string{"msg"})
	if !g.Match(`{"level":"error","msg":"boom happened"}`) {
		t.Error("expected match on msg field")
	}
	if g.Match(`{"level":"boom","msg":"ok"}`) {
		t.Error("expected no match when pattern only in other field")
	}
}

func TestGrep_Match_CaseInsensitiveKey(t *testing.T) {
	g, _ := NewGrep(`hello`, []string{"Message"})
	if !g.Match(`{"message":"hello world"}`) {
		t.Error("expected case-insensitive key match")
	}
}

func TestGrep_Match_NonJSON_WithFields(t *testing.T) {
	g, _ := NewGrep(`error`, []string{"msg"})
	if g.Match(`not json at all`) {
		t.Error("expected no match for non-JSON when fields specified")
	}
}

func TestGrep_Match_NumericField(t *testing.T) {
	g, _ := NewGrep(`42`, []string{"code"})
	if !g.Match(`{"code":42}`) {
		t.Error("expected match on numeric field")
	}
}

func TestGrep_Match_MultipleFields(t *testing.T) {
	g, _ := NewGrep(`target`, []string{"msg", "service"})
	if !g.Match(`{"msg":"ok","service":"target-svc"}`) {
		t.Error("expected match when pattern in second field")
	}
}
