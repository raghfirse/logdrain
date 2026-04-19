package formatter

import (
	"strings"
	"testing"
)

func TestLevelColor_KnownLevels(t *testing.T) {
	cases := []struct {
		level string
		want  string
	}{
		{"debug", colorCyan},
		{"DEBUG", colorCyan},
		{"info", colorGreen},
		{"INFO", colorGreen},
		{"warn", colorYellow},
		{"WARNING", colorYellow},
		{"error", colorRed},
		{"FATAL", colorRed},
		{"unknown", colorGray},
		{"", colorGray},
	}
	for _, tc := range cases {
		got := LevelColor(tc.level)
		if got != tc.want {
			t.Errorf("LevelColor(%q) = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestColorize_WrapsText(t *testing.T) {
	result := Colorize(colorRed, "hello")
	if !strings.Contains(result, "hello") {
		t.Error("expected colorized string to contain original text")
	}
	if !strings.HasPrefix(result, colorRed) {
		t.Error("expected colorized string to start with color code")
	}
	if !strings.HasSuffix(result, colorReset) {
		t.Error("expected colorized string to end with reset code")
	}
}

func TestSourceColor_Deterministic(t *testing.T) {
	name := "my-service"
	if SourceColor(name) != SourceColor(name) {
		t.Error("SourceColor should be deterministic for the same name")
	}
}

func TestSourceColor_DifferentNames(t *testing.T) {
	// Not guaranteed to differ, but a basic smoke test.
	_ = SourceColor("service-a")
	_ = SourceColor("service-b")
	_ = SourceColor("")
}
