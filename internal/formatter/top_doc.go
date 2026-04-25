// Package formatter provides log line transformation and analysis primitives.
//
// # Top
//
// The Top processor tracks the most frequent values for a specified JSON field
// across all log lines and emits a ranked summary at flush time.
//
// Usage:
//
//	--top level:5
//
// This suppresses individual log lines and instead accumulates value counts.
// On program exit (or explicit Flush), it emits up to N JSON lines of the form:
//
//	{"level":"error","count":42}
//	{"level":"warn","count":17}
//	{"level":"info","count":5}
//
// Lines that are not valid JSON are passed through unchanged.
// If the specified field is absent from a line, that line is silently skipped.
package formatter
