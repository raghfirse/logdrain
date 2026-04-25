// Package formatter provides log formatting, filtering, and transformation
// primitives used by logdrain.
//
// # Pick
//
// The Pick transformer retains only a specified set of top-level JSON fields
// from each log line, discarding all others. This is useful for reducing
// noise when you only care about a small subset of structured fields.
//
// Non-JSON lines are passed through unchanged.
//
// Usage via CLI flag:
//
//	--pick msg,level,ts
//
// Keys are matched case-insensitively, so --pick MSG will match a field
// named "msg" in the JSON object.
package formatter
