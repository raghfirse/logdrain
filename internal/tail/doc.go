// Package tail provides utilities for tailing log files, including
// support for log rotation detection. It exposes two main types:
//
//   - Tailer: seeks to the end of a file on open and emits new lines
//     as they are appended, using polling.
//
//   - RotatingTailer: extends Tailer behaviour with rotation detection.
//     When a file shrinks (truncation or replacement), the file is
//     reopened automatically so no lines are missed after rotation.
//
// Both types emit lines over a channel and respect context cancellation.
package tail
