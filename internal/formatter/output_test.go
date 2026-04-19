package formatter

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestMultiWriter_WritesToAll(t *testing.T) {
	var a, b bytes.Buffer
	mw := NewMultiWriter(&a, &b)
	_, err := mw.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.String() != "hello" || b.String() != "hello" {
		t.Errorf("expected both writers to receive 'hello', got %q and %q", a.String(), b.String())
	}
}

func TestMultiWriter_Add(t *testing.T) {
	var a, b bytes.Buffer
	mw := NewMultiWriter(&a)
	mw.Add(&b)
	mw.Write([]byte("x"))
	if b.String() != "x" {
		t.Errorf("expected 'x', got %q", b.String())
	}
}

func TestMultiWriter_ReturnsFirstError(t *testing.T) {
	bad := &errWriter{err: errors.New("write failed")}
	var good bytes.Buffer
	mw := NewMultiWriter(bad, &good)
	_, err := mw.Write([]byte("data"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCountingWriter_CountsBytes(t *testing.T) {
	var buf bytes.Buffer
	cw := NewCountingWriter(&buf)
	cw.Write([]byte("hello"))
	cw.Write([]byte(" world"))
	if cw.BytesWritten() != 11 {
		t.Errorf("expected 11, got %d", cw.BytesWritten())
	}
	if buf.String() != "hello world" {
		t.Errorf("unexpected buf content: %q", buf.String())
	}
}

func TestCountingWriter_ZeroInitial(t *testing.T) {
	cw := NewCountingWriter(io.Discard)
	if cw.BytesWritten() != 0 {
		t.Errorf("expected 0, got %d", cw.BytesWritten())
	}
}

type errWriter struct{ err error }

func (e *errWriter) Write(_ []byte) (int, error) { return 0, e.err }
