package formatter

import "sync/atomic"

// Sampler drops all but every Nth log line per key.
type Sampler struct {
	n       uint64
	counts  map[string]*uint64
}

// NewSampler creates a Sampler that emits 1 out of every n lines.
// A value of 0 or 1 disables sampling (all lines pass).
func NewSampler(n uint64) *Sampler {
	return &Sampler{n: n, counts: make(map[string]*uint64)}
}

// Sample returns true if the line should be emitted.
func (s *Sampler) Sample(key string) bool {
	if s.n <= 1 {
		return true
	}
	ctr, ok := s.counts[key]
	if !ok {
		var v uint64
		s.counts[key] = &v
		ctr = &v
	}
	v := atomic.AddUint64(ctr, 1)
	return v%s.n == 1
}

// Reset clears all counters.
func (s *Sampler) Reset() {
	s.counts = make(map[string]*uint64)
}
