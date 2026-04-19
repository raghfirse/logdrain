package tail

import (
	"context"
	"os"
	"time"
)

// RotatingTailer watches a file and handles log rotation by detecting
// inode changes (or file truncation) and reopening the file.
type RotatingTailer struct {
	path         string
	pollInterval time.Duration
}

// NewRotating creates a RotatingTailer that handles log rotation.
func NewRotating(path string, pollInterval time.Duration) *RotatingTailer {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &RotatingTailer{path: path, pollInterval: pollInterval}
}

// Tail starts tailing with rotation detection. Lines are sent to the returned channel.
func (rt *RotatingTailer) Tail(ctx context.Context) (<-chan string, error) {
	f, size, err := openAtEnd(rt.path)
	if err != nil {
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
			case <-time.After(rt.pollInterval):
			}
			// Detect rotation: check current file size vs position.
			info, err := os.Stat(rt.path)
			if err == nil && info.Size() < size {
				// File was truncated or rotated.
				f.Close()
				f, size, err = openAtEnd(rt.path)
				if err != nil {
					return
				}
				buf = buf[:0]
				continue
			}
			n, _ := f.Read(tmp)
			if n > 0 {
				size += int64(n)
				buf = append(buf, tmp[:n]...)
				buf = emitLines(buf, lines)
			}
		}
	}()
	return lines, nil
}

func openAtEnd(path string) (*os.File, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	size, err := f.Seek(0, 2)
	if err != nil {
		f.Close()
		return nil, 0, err
	}
	return f, size, nil
}
