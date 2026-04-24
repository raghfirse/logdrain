// Package formatter provides log line formatting, filtering, and transformation.
//
// # Grep
//
// The Grep filter matches log lines against a regular expression. It can
// operate on the entire raw line or be scoped to specific JSON fields.
//
// # Flag syntax
//
//	--grep "pattern"               # match anywhere in the line
//	--grep "pattern:field"         # match within a single JSON field
//	--grep "pattern:field1,field2" # match within multiple JSON fields
//
// Patterns are standard Go regular expressions (RE2 syntax). Matching is
// case-sensitive by default; use (?i) for case-insensitive matching.
//
// When fields are specified, only lines that contain valid JSON with at least
// one of the named fields are considered. Lines that do not parse as JSON, or
// that lack all specified fields, are dropped unless --grep-invert is set.
//
// # Examples
//
//	--grep "timeout"               # any line containing "timeout"
//	--grep "5[0-9]{2}:status"      # lines where status field is 5xx
//	--grep "db.*error:msg,service" # lines where msg or service matches
//	--grep "(?i)warn"              # case-insensitive match for "warn"
package formatter
