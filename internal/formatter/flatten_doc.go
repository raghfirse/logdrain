// Package formatter provides formatting, filtering, and transformation
// utilities for structured JSON log lines.
//
// # Flatten
//
// The Flattener collapses nested JSON objects into a single-level map
// using dot-notation (or a custom separator) for compound keys.
//
// Example input:
//
//	{"request":{"method":"GET","path":"/api"},"status":200}
//
// Example output (separator "."):
//
//	{"request.method":"GET","request.path":"/api","status":200}
//
// Use --flatten=. (or --flatten=_ etc.) on the CLI to enable.
package formatter
