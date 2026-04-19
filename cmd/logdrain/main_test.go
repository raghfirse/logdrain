package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logdrain/internal/formatter"
	"github.com/user/logdrain/internal/filter"
)

func TestFilterExprParsing_Valid(t *testing.T) {
	expressions := []string{"level=error", "service=api", "level=info"}
	for _, expr := range expressions {
		_, err := filter.New(expr)
		if err != nil {
			t.Errorf("expected valid expression %q, got error: %v", expr, err)
		}
	}
}

func TestFilterExprParsing_Invalid(t *testing.T) {
	_, err := filter.New("notression")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestFormatterIntegration_Raw(t *testing.T) {
	var buf bytes.Buffer
	fmt := formatter.New(&buf, formatter.Raw)
	line := `{"level":"info","msg":"hello"}`
	if err := fmt.Write("src", line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), line) {
		t.Errorf("expected output to contain %q, got %q", line, buf.String())
	}
}

func TestFormatterIntegration_Pretty(t *testing.T) {
	var buf bytes.Buffer
	fmt := formatter.New(&buf, formatter.Pretty)
	line := `{"level":"warn","msg":"something happened"}`
	if err := fmt.Write("myapp", line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "warn") {
		t.Errorf("expected output to contain level, got: %q", out)
	}
	if !strings.Contains(out, "something happened") {
		t.Errorf("expected output to contain message, got: %q", out)
	}
}
