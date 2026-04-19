package formatter

import "fmt"

// ANSI color codes for pretty output.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// LevelColor returns an ANSI color code for a given log level string.
func LevelColor(level string) string {
	switch level {
	case "debug", "DEBUG":
		return colorCyan
	case "info", "INFO":
		return colorGreen
	case "warn", "WARN", "warning", "WARNING":
		return colorYellow
	case "error", "ERROR", "fatal", "FATAL":
		return colorRed
	default:
		return colorGray
	}
}

// Colorize wraps text in the given ANSI color code.
func Colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, colorReset)
}

// SourceColor returns a stable color for a named source using a small palette.
func SourceColor(name string) string {
	palette := []string{colorBlue, colorCyan, colorYellow, colorGreen, colorRed}
	h := 0
	for _, c := range name {
		h = (h*31 + int(c)) % len(palette)
	}
	return palette[h]
}
