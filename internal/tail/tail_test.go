package tail

import (
	"context"
	"os"
	"testing"
	"time"
)

func writeTmp(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "tail_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	if content != "" {
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
	}
	return f
}

func TestTail_EmitsNewLines(t *testing.T) {
	f := writeTmp(t, "")
	defer os.Remove(f.Name())
	defer f.Close()

	tr := New(f.Name(), 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines, err := tr.Tail(ctx)
	if err != nil {
		t.Fatal(err)
	}

	f.WriteString("hello\nworld\n")

	got := []string{}
	for line := range lines {
		got = append(got, line)
		if len(got) == 2 {
			cancel()
		}
	}
	if len(got) != 2 || got[0] != "hello" || got[1] != "world" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestTail_IgnoresExistingContent(t *testing.T) {
	f := writeTmp(t, "existing\n")
	defer os.Remove(f.Name())
	defer f.Close()

	tr := New(f.Name(), 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	lines, err := tr.Tail(ctx)
	if err != nil {
		t.Fatal(err)
	}

	got := []string{}
	for line := range lines {
		got = append(got, line)
	}
	if len(got) != 0 {
		t.Errorf("expected no lines, got: %v", got)
	}
}

func TestTail_InvalidPath(t *testing.T) {
	tr := New("/nonexistent/path/file.log", 0)
	_, err := tr.Tail(context.Background())
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestEmitLines(t *testing.T) {
	ch := make(chan string, 10)
	remaining := emitLines([]byte("foo\nbar\nbaz"), ch)
	close(ch)
	lines := []string{}
	for l := range ch {
		lines = append(lines, l)
	}
	if len(lines) != 2 || lines[0] != "foo" || lines[1] != "bar" {
		t.Errorf("unexpected: %v", lines)
	}
	if string(remaining) != "baz" {
		t.Errorf("unexpected remainder: %q", remaining)
	}
}
