package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// AggregateMode defines how values are combined.
type AggregateMode int

const (
	AggregateCount AggregateMode = iota
	AggregateSum
	AggregateMin
	AggregateMax
)

// Aggregator accumulates numeric values per key from JSON log lines.
type Aggregator struct {
	mu     sync.Mutex
	field  string
	mode   AggregateMode
	bucket map[string]float64
	counts map[string]int
}

// NewAggregator creates an Aggregator for the given field and mode.
func NewAggregator(field string, mode AggregateMode) *Aggregator {
	return &Aggregator{
		field:  field,
		mode:   mode,
		bucket: make(map[string]float64),
		counts: make(map[string]int),
	}
}

// Add ingests a JSON log line, extracting the target field value.
func (a *Aggregator) Add(line string) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return
	}
	v, ok := obj[a.field]
	if !ok {
		return
	}
	var num float64
	switch val := v.(type) {
	case float64:
		num = val
	default:
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	key := a.field
	a.counts[key]++
	switch a.mode {
	case AggregateCount:
		a.bucket[key]++
	case AggregateSum:
		a.bucket[key] += num
	case AggregateMin:
		if a.counts[key] == 1 || num < a.bucket[key] {
			a.bucket[key] = num
		}
	case AggregateMax:
		if a.counts[key] == 1 || num > a.bucket[key] {
			a.bucket[key] = num
		}
	}
}

// Report returns a sorted slice of "key=value" summary strings.
func (a *Aggregator) Report() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	keys := make([]string, 0, len(a.bucket))
	for k := range a.bucket {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s=%.4g", k, a.bucket[k]))
	}
	return out
}

// Reset clears all accumulated state.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.bucket = make(map[string]float64)
	a.counts = make(map[string]int)
}
