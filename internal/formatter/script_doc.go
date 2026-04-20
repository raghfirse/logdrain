// Package formatter — script module
//
// The script feature lets users add or overwrite fields on every JSON log line
// using Go text/template expressions.
//
// Usage (CLI flag, repeatable):
//
//	--script 'env=production'
//	--script 'summary={{.level}}: {{.msg}}'
//
// Rules are evaluated in declaration order. Each rule has the form:
//
//	key=template
//
// Where:
//   - key   — the JSON field name to set (created or overwritten)
//   - template — a Go text/template string with access to all top-level
//     fields of the current log line as {{ .fieldname }}
//
// Non-JSON lines pass through unchanged.
package formatter
