// Package formatter provides utilities for formatting, coloring, and writing
// structured JSON log entries to one or more output destinations.
//
// Output targets are abstracted as Sink values, which wrap any io.Writer and
// optionally handle closing file-backed destinations.
//
// MultiWriter fans writes out to several sinks simultaneously, useful when
// output should go to both stdout and a log file. CountingWriter tracks how
// many bytes have been emitted, enabling lightweight throughput metrics.
//
// Typical usage:
//
//	// Write to stdout and a log file simultaneously.
//	 fileSink, err := formatter.NewFileSink("/var/log/app.log")
//	 if err != nil {
//	     log.Fatal(err)
//	 }
//	 defer fileSink.Close()
//
//	 mw := formatter.NewMultiWriter(formatter.StdoutSink(), fileSink)
//	 mw.Write(entry)
package formatter
