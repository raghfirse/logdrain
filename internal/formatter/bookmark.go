package formatter

import (
	"encoding/json"
	"os"
	"sync"
)

// Bookmark tracks the last-seen log line per source, persisted to a JSON file.
// This allows logdrain to resume from where it left off across restarts.
type Bookmark struct {
	mu       sync.Mutex
	path     string
	positions map[string]int64
}

// NewBookmark loads an existing bookmark file or returns an empty one.
func NewBookmark(path string) (*Bookmark, error) {
	b := &Bookmark{
		path:      path,
		positions: make(map[string]int64),
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &b.positions); err != nil {
		return nil, err
	}
	return b, nil
}

// Get returns the last recorded offset for the given source name.
func (b *Bookmark) Get(source string) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.positions[source]
}

// Set updates the offset for the given source name and flushes to disk.
func (b *Bookmark) Set(source string, offset int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.positions[source] = offset
	return b.flush()
}

// flush writes the current positions map to disk. Caller must hold mu.
func (b *Bookmark) flush() error {
	data, err := json.MarshalIndent(b.positions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Reset clears all bookmarks and removes the file.
func (b *Bookmark) Reset() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.positions = make(map[string]int64)
	err := os.Remove(b.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
