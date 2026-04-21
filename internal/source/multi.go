package source

import (
	"context"
	"sync"
)

// Entry is a line paired with the name of its originating source.
type Entry struct {
	Source string
	Line   string
}

// Merge fans-in lines from multiple sources into a single channel.
// Each entry is tagged with the source name. The returned channel is
// closed once all sources are exhausted or ctx is cancelled.
func Merge(ctx context.Context, sources []*Source) <-chan Entry {
	out := make(chan Entry)
	var wg sync.WaitGroup
	for _, s := range sources {
		wg.Add(1)
		go func(s *Source) {
			defer wg.Done()
			for line := range s.Lines(ctx) {
				select {
				case <-ctx.Done():
					return
				case out <- Entry{Source: s.Name, Line: line}:
				}
			}
		}(s)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Filter returns a new channel containing only entries whose Source name
// matches one of the provided names. If names is empty, all entries are
// passed through. The returned channel is closed when in is closed.
func Filter(ctx context.Context, in <-chan Entry, names ...string) <-chan Entry {
	out := make(chan Entry)
	allowed := make(map[string]struct{}, len(names))
	for _, n := range names {
		allowed[n] = struct{}{}
	}
	go func() {
		defer close(out)
		for entry := range in {
			if len(allowed) > 0 {
				if _, ok := allowed[entry.Source]; !ok {
					continue
				}
			}
			select {
			case <-ctx.Done():
				return
			case out <- entry:
			}
		}
	}()
	return out
}
