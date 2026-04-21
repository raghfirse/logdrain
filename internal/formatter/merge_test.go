package formatter

import (
	"encoding/json"
	"testing"
)

func TestNewMerge_Empty(t *testing.T) {
	m, err := NewMerge("", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil merger")
	}
}

func TestNewMerge_InvalidJSON(t *testing.T) {
	_, err := NewMerge("not-json", false)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestMerge_Apply_NonJSON(t *testing.T) {
	m, _ := NewMerge(`{"env":"prod"}`, false)
	result := m.Apply("plain text line")
	if result != "plain text line" {
		t.Errorf("expected passthrough, got %q", result)
	}
}

func TestMerge_Apply_AddsFields(t *testing.T) {
	m, _ := NewMerge(`{"env":"prod","region":"us-east-1"}`, false)
	result := m.Apply(`{"msg":"hello"}`)
	var out map[string]string
	if err := json.Unmarshal([]byte(result), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", out["env"])
	}
	if out["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", out["region"])
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %q", out["msg"])
	}
}

func TestMerge_Apply_NoOverwrite(t *testing.T) {
	m, _ := NewMerge(`{"env":"prod"}`, false)
	result := m.Apply(`{"env":"dev","msg":"hi"}`)
	var out map[string]string
	json.Unmarshal([]byte(result), &out)
	if out["env"] != "dev" {
		t.Errorf("expected env=dev (no overwrite), got %q", out["env"])
	}
}

func TestMerge_Apply_Overwrite(t *testing.T) {
	m, _ := NewMerge(`{"env":"prod"}`, true)
	result := m.Apply(`{"env":"dev","msg":"hi"}`)
	var out map[string]string
	json.Unmarshal([]byte(result), &out)
	if out["env"] != "prod" {
		t.Errorf("expected env=prod (overwrite), got %q", out["env"])
	}
}

func TestMerge_Apply_EmptyFields(t *testing.T) {
	m, _ := NewMerge("", false)
	line := `{"msg":"unchanged"}`
	result := m.Apply(line)
	if result != line {
		t.Errorf("expected unchanged line, got %q", result)
	}
}

func TestParseMergeFlag_Valid(t *testing.T) {
	m, err := ParseMergeFlag(`{"app":"logdrain"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := m.Apply(`{"msg":"test"}`)
	var out map[string]string
	json.Unmarshal([]byte(result), &out)
	if out["app"] != "logdrain" {
		t.Errorf("expected app=logdrain, got %q", out["app"])
	}
}
