// Package formatter provides log line formatting, filtering, and output utilities.
//
// # Deduplication
//
// The DedupeFilter suppresses repeated identical log lines within a configurable
// time window. This is useful when tailing noisy services that emit the same
// error repeatedly.
//
// Usage:
//
//	df := formatter.NewDedupeFilter(5 * time.Second)
//	if !df.IsDuplicate(line) {
//	    // emit line
//	}
//
// The DedupeFlag helper integrates with the CLI flag system.
package formatter
