package formatter

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func makeLineWithAfterTime(t time.Time) string {
	b, _ := json.Marshal(map[string]string{"time": t.Format(time.RFC3339), "msg": "hello"})
	return string(b)
}

func TestAfter_PassesOlderLine(t *testing.T) {
	cutoff := time.Now()
	a, _ := NewAfter(cutoff)
	old := makeLineWithAfterTime(cutoff.Add(-time.Minute))
	out, ok := a.Apply(old)
	if !ok {
		t.Fatal("expected older line to pass through")
	}
	if out != old {
		t.Fatalf("expected line unchanged, got %q", out)
	}
}

func TestAfter_SuppressesNewerLine(t *testing.T) {
	cutoff := time.Now()
	a, _ := NewAfter(cutoff)
	newer := makeLineWithAfterTime(cutoff.Add(time.Minute))
	_, ok := a.Apply(newer)
	if ok {
		t.Fatal("expected newer line to be suppressed")
	}
}

func TestAfter_PassesExactCutoff(t *testing.T) {
	cutoff, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	a, _ := NewAfter(cutoff)
	exact := makeLineWithAfterTime(cutoff)
	_, ok := a.Apply(exact)
	if !ok {
		t.Fatal("expected line at exact cutoff to pass through")
	}
}

func TestAfter_NonJSONPassthrough(t *testing.T) {
	a, _ := NewAfter(time.Now())
	line := "not json at all"
	out, ok := a.Apply(line)
	if !ok {
		t.Fatal("expected non-JSON to pass through")
	}
	if out != line {
		t.Fatalf("expected line unchanged")
	}
}

func TestAfter_MissingTimestampPassthrough(t *testing.T) {
	a, _ := NewAfter(time.Now())
	line := `{"msg":"no timestamp here"}`
	_, ok := a.Apply(line)
	if !ok {
		t.Fatal("expected line without timestamp to pass through")
	}
}

func TestNewAfter_ZeroCutoffError(t *testing.T) {
	_, err := NewAfter(time.Time{})
	if err == nil {
		t.Fatal("expected error for zero cutoff")
	}
}

func TestParseAfterFlag_RFC3339(t *testing.T) {
	a, err := ParseAfterFlag("2024-06-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil After")
	}
}

func TestParseAfterFlag_NegativeDuration(t *testing.T) {
	a, err := ParseAfterFlag("-2h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil After")
	}
}

func TestParseAfterFlag_Empty(t *testing.T) {
	_, err := ParseAfterFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}

func TestParseAfterFlag_Invalid(t *testing.T) {
	_, err := ParseAfterFlag("not-a-time")
	if err == nil {
		t.Fatal("expected error for invalid value")
	}
}

func TestAfter_UnixTimestamp(t *testing.T) {
	cutoff := time.Unix(1700000000, 0)
	a, _ := NewAfter(cutoff)
	older := fmt.Sprintf(`{"ts":%d,"msg":"old"}`, cutoff.Unix()-60)
	_, ok := a.Apply(older)
	if !ok {
		t.Fatal("expected older unix timestamp line to pass through")
	}
	newer := fmt.Sprintf(`{"ts":%d,"msg":"new"}`, cutoff.Unix()+60)
	_, ok = a.Apply(newer)
	if ok {
		t.Fatal("expected newer unix timestamp line to be suppressed")
	}
}
