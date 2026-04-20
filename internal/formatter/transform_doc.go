// Package formatter — transform module
//
// The transform feature allows renaming JSON fields on the fly as log lines
// are processed. This is useful when different log sources use different
// field names for the same concept (e.g. "msg" vs "message", "ts" vs "time").
//
// Usage via CLI flag:
//
//	--transform msg:message --transform ts:timestamp
//
// Rules are applied in order. Key matching is case-insensitive. If the source
// key does not exist in a given line the rule is silently skipped. Non-JSON
// lines are passed through unchanged.
package formatter
