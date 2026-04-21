// Package formatter provides log line transformation and filtering primitives.
//
// # Since Filter
//
// The Since filter discards log lines whose timestamp predates a given cutoff.
// It is intended for use with --since <duration>, e.g.:
//
//	--since 5m   # only show lines from the last 5 minutes
//	--since 1h   # only show lines from the last hour
//
// Timestamp detection
//
// The filter inspects the field named by --since-field (default: "time").
// If that field is absent it falls back to common aliases in order:
// "ts", "timestamp", "@timestamp", "t".
//
// Lines that are not valid JSON, or whose timestamp cannot be parsed,
// are passed through unchanged rather than silently dropped.
package formatter
