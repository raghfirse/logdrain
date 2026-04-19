package tail

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRotatingTailer_EmitsNewLines(t *testing.T) {
	f := writeTmp(t, "")
	defer os.Remove(f.Name())
	defer f.Close()

	rt := NewRotating(f.Name(), 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines, err := rt.Tail(ctx)
	if err != nil {
		t.Fatal(err)
	}

	f.WriteString("line1\nline2\n")

	got := []string{}
	for line := range lines {
		got = append(got, line)
		if len(got) == 2 {
			cancel()
		}
	}
	if len(got) != 2 || got[0] != "line1" || got[1] != "line2" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestRotatingTailer_HandlesRotation(t *testing.T) {
	f := writeTmp(t, "old content\n")
	name := f.Name()
	defer os.Remove(name)

	rt := NewRotating(name, 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines, err := rt.Tail(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate rotation by truncating.
	time.Sleep(30 * time.Millisecond)
	os.WriteFile(name, []byte("new\n"), 0644)

	got := []string{}
	for line := range lines {
		got = append(got, line)
		if len(got) >= 1 {
			cancel()
		}
	}
	if len(got) == 0 {
		t.Error("expected at least one line after rotation")
	}
}

func TestNewRotating_InvalidPath(t *testing.T) {
	rt := NewRotating("/no/such/file.log", 0)
	_, err := rt.Tail(context.Background())
	if err == nil {
		t.Error("expected error")
	}
}
