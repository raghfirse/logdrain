package formatter

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSink_Write(t *testing.T) {
	var buf bytes.Buffer
	s := NewSink(&buf)
	s.Write([]byte("test line\n"))
	if buf.String() != "test line\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestNewSink_CloseNoop(t *testing.T) {
	var buf bytes.Buffer
	s := NewSink(&buf)
	if err := s.Close(); err != nil {
		t.Errorf("expected nil error on Close, got %v", err)
	}
}

func TestNewFileSink_WritesAndCloses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")
	s, err := NewFileSink(path)
	if err != nil {
		t.Fatalf("NewFileSink: %v", err)
	}
	s.Write([]byte("hello file\n"))
	s.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "hello file\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestNewFileSink_InvalidPath(t *testing.T) {
	_, err := NewFileSink("/nonexistent/dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestStdoutSink_NotNil(t *testing.T) {
	if StdoutSink() == nil {
		t.Error("StdoutSink() returned nil")
	}
}

func TestStderrSink_NotNil(t *testing.T) {
	if StderrSink() == nil {
		t.Error("StderrSink() returned nil")
	}
}
