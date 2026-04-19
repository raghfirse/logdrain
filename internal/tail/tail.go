package tail

import (
	"context"
	"io"
	"os"
	"time"
)

// Tailer watches a file and emits new lines as they are appended.
type Tailer struct {
	path     string
	pollInterval time.Duration
}

// New creates a new Tailer for the given file path.
func New(path string, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &Tailer{path: path, pollInterval: pollInterval}
}

// Tail opens the file, seeks to the end, and emits new lines to the returned channel.
// The channel is closed when ctx is cancelled or an error occurs.
func (t *Tailer) Tail(ctx context.Context) (<-chan string, error) {
	f, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}

	// Seek to end so we only emit new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}

	lines := make(chan string, 64)
	go func() {
		defer close(lines)
		defer f.Close()
		buf := make([]byte, 0, 4096)
		tmp := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, err := f.Read(tif n > 0 {
			uf = append(buf, tmp[:n]...uf = emitLines(buf, lines)
			}
			if err == io.EOF {
				select {
				case <-ctx.Done():
					return
				case <-time.After(t.pollInterval):
				}
				continue
			}
			if err != nil {
				return
			}
		}
	}()
	return lines, nil
}

// emitLines sends complete newline-terminated lines to ch and returns remaining bytes.
func emitLines(buf []byte, ch chan<- string) []byte {
	start := 0
	for i, b := range buf {
		if b == '\n' {
			line := string(buf[start:i])
			ch <- line
			start = i + 1
		}
	}
	return buf[start:]
}
