package source

import (
	"bufio"
	"context"
	"io"
	"os"
)

// Source represents a log source that emits lines.
type Source struct {
	Name   string
	reader io.ReadCloser
}

// New creates a Source from a named reader.
func New(name string, r io.ReadCloser) *Source {
	return &Source{Name: name, reader: r}
}

// NewFromFile opens a file and returns a Source.
func NewFromFile(path string) (*Source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return New(path, f), nil
}

// NewFromStdin returns a Source reading from os.Stdin.
func NewFromStdin() *Source {
	return New("stdin", io.NopCloser(os.Stdin))
}

// Lines streams lines from the source into the returned channel.
// The channel is closed when EOF is reached or ctx is cancelled.
func (s *Source) Lines(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer s.reader.Close()
		scanner := bufio.NewScanner(s.reader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		}
	}()
	return ch
}
