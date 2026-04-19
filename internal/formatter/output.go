package formatter

import (
	"io"
	"sync"
)

// MultiWriter fans out writes to multiple io.Writer targets.
type MultiWriter struct {
	mu      sync.Mutex
	writers []io.Writer
}

// NewMultiWriter creates a MultiWriter that writes to all provided writers.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Add appends a writer to the fan-out list.
func (m *MultiWriter) Add(w io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.writers = append(m.writers, w)
}

// Write writes p to all underlying writers. Returns the first error encountered.
func (m *MultiWriter) Write(p []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, w := range m.writers {
		if _, err := w.Write(p); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

// CountingWriter wraps an io.Writer and counts bytes written.
type CountingWriter struct {
	mu    sync.Mutex
	inner io.Writer
	total int64
}

// NewCountingWriter wraps w with byte counting.
func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{inner: w}
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.inner.Write(p)
	c.mu.Lock()
	c.total += int64(n)
	c.mu.Unlock()
	return n, err
}

// BytesWritten returns the total number of bytes written so far.
func (c *CountingWriter) BytesWritten() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total
}
