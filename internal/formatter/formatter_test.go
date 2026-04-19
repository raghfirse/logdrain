package formatter

import (
	"bytes"
	"strings"
	"testing"
)

func TestWrite_RawFormat(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatRaw)
	if err := f.Write("app", `{"level":"info","msg":"hello"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[app]") {
		t.Errorf("expected source prefix, got: %s", got)
	}
	if !strings.Contains(got, `{"level":"info"`) {
		t.Errorf("expected raw JSON in output, got: %s", got)
	}
}

func TestWrite_PrettyFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatPretty)
	err := f.Write("svc", `{"level":"error","msg":"something failed","time":"2024-01-02T15:04:05Z"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[svc]") {
		t.Errorf("expected source name, got: %s", got)
	}
	if !strings.Contains(got, "ERROR") {
		t.Errorf("expected level ERROR, got: %s", got)
	}
	if !strings.Contains(got, "something failed") {
		t.Errorf("expected message, got: %s", got)
	}
	if !strings.Contains(got, "15:04:05") {
		t.Errorf("expected formatted time, got: %s", got)
	}
}

func TestWrite_PrettyFormat_InvalidJSON(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatPretty)
	plain := "not json at all"
	if err := f.Write("src", plain); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, plain) {
		t.Errorf("expected raw line in fallback output, got: %s", got)
	}
}

func TestWrite_PrettyFormat_MessageAlias(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatPretty)
	if err := f.Write("x", `{"level":"warn","message":"alt field"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alt field") {
		t.Errorf("expected message from 'message' field, got: %s", buf.String())
	}
}
