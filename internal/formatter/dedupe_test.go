package formatter

import (
	"testing"
	"time"
)

func TestDedupeFilter_FirstOccurrenceNotDuplicate(t *testing.T) {
	d := NewDedupeFilter(5 * time.Second)
	if d.IsDuplicate("hello") {
		t.Fatal("expected first occurrence to not be a duplicate")
	}
}

func TestDedupeFilter_SecondOccurrenceIsDuplicate(t *testing.T) {
	d := NewDedupeFilter(5 * time.Second)
	d.IsDuplicate("hello")
	if !d.IsDuplicate("hello") {
		t.Fatal("expected second occurrence to be a duplicate")
	}
}

func TestDedupeFilter_DifferentLinesNotDuplicate(t *testing.T) {
	d := NewDedupeFilter(5 * time.Second)
	d.IsDuplicate("line1")
	if d.IsDuplicate("line2") {
		t.Fatal("different line should not be a duplicate")
	}
}

func TestDedupeFilter_ExpiresAfterWindow(t *testing.T) {
	now := time.Now()
	d := NewDedupeFilter(2 * time.Second)
	d.now = func() time.Time { return now }
	d.IsDuplicate("msg")

	// advance past window
	d.now = func() time.Time { return now.Add(3 * time.Second) }
	if d.IsDuplicate("msg") {
		t.Fatal("entry should have expired")
	}
}

func TestDedupeFilter_Reset(t *testing.T) {
	d := NewDedupeFilter(5 * time.Second)
	d.IsDuplicate("msg")
	d.Reset()
	if d.IsDuplicate("msg") {
		t.Fatal("after reset, line should not be duplicate")
	}
}

func TestDedupeFilter_ZeroWindow(t *testing.T) {
	now := time.Now()
	d := NewDedupeFilter(0)
	d.now = func() time.Time { return now }
	d.IsDuplicate("msg")
	// with zero window everything expires immediately on next call
	if d.IsDuplicate("msg") {
		t.Fatal("zero window: should not be duplicate on next call")
	}
}
