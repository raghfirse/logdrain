package formatter

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewWindow_Valid(t *testing.T) {
	w, err := NewWindow("level", 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestNewWindow_EmptyField(t *testing.T) {
	_, err := NewWindow("", 10*time.Second)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewWindow_ZeroDuration(t *testing.T) {
	_, err := NewWindow("level", 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestWindow_Apply_NoFlushBeforeExpiry(t *testing.T) {
	w, _ := NewWindow("level", 10*time.Second)
	now := time.Now()
	w.now = func() time.Time { return now }

	out := w.Apply(`{"level":"info","msg":"hello"}`)
	if out != "" {
		t.Fatalf("expected empty before window expires, got %q", out)
	}
}

func TestWindow_Apply_FlushesAfterExpiry(t *testing.T) {
	w, _ := NewWindow("level", 5*time.Second)
	base := time.Now()
	calls := 0
	w.now = func() time.Time {
		calls++
		if calls == 1 {
			return base
		}
		return base.Add(6 * time.Second)
	}

	w.Apply(`{"level":"info"}`)
	out := w.Apply(`{"level":"error"}`)
	if out == "" {
		t.Fatal("expected summary line after window expires")
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("summary is not valid JSON: %v", err)
	}
	if m["_field"] != "level" {
		t.Errorf("expected _field=level, got %v", m["_field"])
	}
}

func TestWindow_Flush_ReturnsCountsAndResets(t *testing.T) {
	w, _ := NewWindow("status", 30*time.Second)
	w.Apply(`{"status":"200"}`)
	w.Apply(`{"status":"200"}`)
	w.Apply(`{"status":"500"}`)

	out := w.Flush()
	if out == "" {
		t.Fatal("expected non-empty flush output")
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(out), &m)
	if m["200"].(float64) != 2 {
		t.Errorf("expected 200 count=2, got %v", m["200"])
	}
	if m["500"].(float64) != 1 {
		t.Errorf("expected 500 count=1, got %v", m["500"])
	}
	// Second flush on empty window should be empty.
	if w.Flush() != "" {
		t.Error("expected empty flush on reset window")
	}
}

func TestParseWindowFlag_Valid(t *testing.T) {
	w, err := ParseWindowFlag("level:10s")
	if err != nil || w == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseWindowFlag_PlainSeconds(t *testing.T) {
	w, err := ParseWindowFlag("level:30")
	if err != nil || w == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseWindowFlag_Empty(t *testing.T) {
	w, err := ParseWindowFlag("")
	if err != nil || w != nil {
		t.Fatal("expected nil window and no error for empty flag")
	}
}

func TestParseWindowFlag_Invalid(t *testing.T) {
	if _, err := ParseWindowFlag("noduration"); err == nil {
		t.Fatal("expected error for missing colon")
	}
	if _, err := ParseWindowFlag("level:notaduration"); err == nil {
		t.Fatal("expected error for bad duration")
	}
}
