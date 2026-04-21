package formatter

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func makeLineWithTime(ts time.Time) string {
	b, _ := json.Marshal(map[string]interface{}{
		"time":    ts.Format(time.RFC3339),
		"message": "hello",
	})
	return string(b)
}

func TestSince_PassesNewerLine(t *testing.T) {
	cutoff := time.Now().Add(-10 * time.Minute)
	s := NewSince(cutoff, "time")
	line := makeLineWithTime(time.Now().Add(-5 * time.Minute))
	if got := s.Apply(line); got != line {
		t.Errorf("expected line to pass, got %q", got)
	}
}

func TestSince_SuppressesOlderLine(t *testing.T) {
	cutoff := time.Now().Add(-10 * time.Minute)
	s := NewSince(cutoff, "time")
	line := makeLineWithTime(time.Now().Add(-20 * time.Minute))
	if got := s.Apply(line); got != "" {
		t.Errorf("expected line suppressed, got %q", got)
	}
}

func TestSince_PassesExactCutoff(t *testing.T) {
	cutoff := time.Now().Add(-10 * time.Minute)
	s := NewSince(cutoff, "time")
	line := makeLineWithTime(cutoff)
	if got := s.Apply(line); got != line {
		t.Errorf("expected line at cutoff to pass, got %q", got)
	}
}

func TestSince_NonJSONPassthrough(t *testing.T) {
	s := NewSince(time.Now(), "time")
	line := "not json at all"
	if got := s.Apply(line); got != line {
		t.Errorf("expected non-JSON passthrough, got %q", got)
	}
}

func TestSince_MissingFieldPassthrough(t *testing.T) {
	s := NewSince(time.Now(), "time")
	line := `{"message":"no timestamp here"}`
	if got := s.Apply(line); got != line {
		t.Errorf("expected missing-field passthrough, got %q", got)
	}
}

func TestSince_AliasTimestamp(t *testing.T) {
	cutoff := time.Now().Add(-10 * time.Minute)
	s := NewSince(cutoff, "")
	b, _ := json.Marshal(map[string]interface{}{
		"ts":      time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		"message": "via alias",
	})
	line := string(b)
	if got := s.Apply(line); got != line {
		t.Errorf("expected alias lookup to pass, got %q", got)
	}
}

func TestParseSinceFlag_Valid(t *testing.T) {
	before := time.Now()
	got, err := ParseSinceFlag("5m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := before.Add(-5 * time.Minute)
	if got.Before(expected.Add(-time.Second)) || got.After(before) {
		t.Errorf("cutoff %v not in expected range", got)
	}
}

func TestParseSinceFlag_Empty(t *testing.T) {
	got, err := ParseSinceFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.IsZero() {
		t.Errorf("expected zero time, got %v", got)
	}
}

func TestParseSinceFlag_Invalid(t *testing.T) {
	_, err := ParseSinceFlag("notaduration")
	if err == nil {
		t.Error("expected error for invalid duration")
	}
}

func TestParseSinceFlag_Negative(t *testing.T) {
	_, err := ParseSinceFlag("-5m")
	if err == nil {
		t.Error("expected error for negative duration")
	}
}

var _ = fmt.Sprintf // avoid unused import
