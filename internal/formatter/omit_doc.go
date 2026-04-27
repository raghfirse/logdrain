// Package formatter — omit
//
// The omit transform removes one or more named keys from every JSON log line
// before the line is forwarded to the output sink. Non-JSON lines are passed
// through unchanged.
//
// Key matching is case-insensitive, so --omit Secret will remove both
// "secret" and "Secret" from the output.
//
// Usage
//
//	--omit key[,key...]
//
// Multiple --omit flags are additive; the following two invocations are
// equivalent:
//
//	logdrain --omit password --omit token
//	logdrain --omit password,token
//
// Example
//
// Input:
//
//	{"level":"info","msg":"login","password":"hunter2","user":"alice"}
//
// Output with --omit password:
//
//	{"level":"info","msg":"login","user":"alice"}
package formatter
