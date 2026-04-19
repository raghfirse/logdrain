package formatter

import (
	"io"
	"os"
)

// Sink represents a destination for formatted log output.
type Sink struct {
	writer io.Writer
	closer io.Closer
}

// NewSink wraps an existing writer. Caller is responsible for closing.
func NewSink(w io.Writer) *Sink {
	return &Sink{writer: w}
}

// NewFileSink opens a file for appending and returns a Sink.
func NewFileSink(path string) (*Sink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &Sink{writer: f, closer: f}, nil
}

// Write writes p to the sink's writer.
func (s *Sink) Write(p []byte) (int, error) {
	return s.writer.Write(p)
}

// Close closes the underlying writer if it implements io.Closer.
func (s *Sink) Close() error {
	if s.closer != nil {
		return s.closer.Close()
	}
	return nil
}

// StdoutSink returns a Sink writing to os.Stdout.
func StdoutSink() *Sink {
	return NewSink(os.Stdout)
}

// StderrSink returns a Sink writing to os.Stderr.
func StderrSink() *Sink {
	return NewSink(os.Stderr)
}
