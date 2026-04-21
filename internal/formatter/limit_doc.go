// Package formatter — limit
//
// The Limiter processor stops the output stream after a fixed number of lines
// have been emitted. It is useful when you only need the first N matching
// entries from a potentially infinite log stream.
//
// Usage (CLI flag):
//
//	--limit 50          # stop after 50 lines
//	--limit 0           # unlimited (default)
//
// When the limit is reached, Apply returns done=true. The caller (main loop)
// should treat this as a clean EOF and exit without printing an error.
package formatter
