package formatter

import (
	"testing"
)

func TestFormatTimestamp_RFC3339(t *testing.T) {
	entry := map[string]any{"time": "2024-03-15T10:20:30Z"}
	got := FormatTimestamp(entry, false)
	want := "2024-03-15 10:20:30.000"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatTimestamp_Short(t *testing.T) {
	entry := map[string]any{"ts": "2024-03-15T10:20:30Z"}
	got := FormatTimestamp(entry, true)
	want := "10:20:30"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatTimestamp_UnixFloat(t *testing.T) {
	// 2024-01-01 00:00:00 UTC = 1704067200
	entry := map[string]any{"timestamp": float64(1704067200)}
	got := FormatTimestamp(entry, false)
	if got == "" {
		t.Error("expected non-empty timestamp for unix float")
	}
}

func TestFormatTimestamp_AliasKeys(t *testing.T) {
	for _, key := range []string{"time", "timestamp", "ts", "@timestamp"} {
		entry := map[string]any{key: "2024-03-15T10:20:30Z"}
		got := FormatTimestamp(entry, false)
		if got == "" {
			t.Errorf("key %q: expected non-empty result", key)
		}
	}
}

func TestFormatTimestamp_Missing(t *testing.T) {
	entry := map[string]any{"message": "hello"}
	got := FormatTimestamp(entry, false)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestFormatTimestamp_InvalidString(t *testing.T) {
	entry := map[string]any{"time": "not-a-time"}
	got := FormatTimestamp(entry, false)
	if got != "" {
		t.Errorf("expected empty string for unparseable time, got %q", got)
	}
}

func TestFormatTimestamp_RFC3339Nano(t *testing.T) {
	entry := map[string]any{"time": "2024-03-15T10:20:30.123456789Z"}
	got := FormatTimestamp(entry, false)
	if got == "" {
		t.Error("expected non-empty timestamp for RFC3339Nano")
	}
}
