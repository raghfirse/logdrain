package source

import (
	"context"
	"io"
	"strings"
	"testing"
)

func nopCloser(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

func TestLines_EmitsAllLines(t *testing.T) {
	src := New("test", nopCloser("line1\nline2\nline3"))
	ctx := context.Background()
	var got []string
	for line := range src.Lines(ctx) {
		got = append(got, line)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line1" || got[1] != "line2" || got[2] != "line3" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestLines_EmptyInput(t *testing.T) {
	src := New("empty", nopCloser(""))
	ctx := context.Background()
	var got []string
	for line := range src.Lines(ctx) {
		got = append(got, line)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(got))
	}
}

func TestLines_CancelContext(t *testing.T) {
	// Large input so the goroutine blocks on send.
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("line\n")
	}
	src := New("cancel", nopCloser(sb.String()))
	ctx, cancel := context.WithCancel(context.Background())
	ch := src.Lines(ctx)
	// Read one line then cancel.
	<-ch
	cancel()
	// Drain remaining; channel must eventually close.
	for range ch {
	}
}

func TestNew_Name(t *testing.T) {
	src := New("myname", nopCloser(""))
	if src.Name != "myname" {
		t.Errorf("expected name 'myname', got %q", src.Name)
	}
}
