// Package formatter provides utilities for formatting, coloring, and writing
// structured JSON log entries to one or more output destinations.
//
// Output targets are abstracted as Sink values, which wrap any io.Writer and
// optionally handle closing file-backed destinations.
//
// MultiWriter fans writes out to several sinks simultaneously, useful when
// output should go to both stdout and a log file. CountingWriter tracks how
// many bytes have been emitted, enabling lightweight throughput metrics.
package formatter
