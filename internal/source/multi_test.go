package source

import (
	"context"
	"sort"
	"testing"
)

func TestMerge_MultipleSourcesAllLines(t *testing.T) {
	a := New("a", nopCloser("a1\na2"))
	b := New("b", nopCloser("b1\nb2"))
	ctx := context.Background()
	var entries []Entry
	for e := range Merge(ctx, []*Source{a, b}) {
		entries = append(entries, e)
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", entries)
	}
}

func TestMerge_SourceNamesPreserved(t *testing.T) {
	a := New("src-a", nopCloser("hello"))
	b := New("src-b", nopCloser("world"))
	ctx := context.Background()
	names := map[string]bool{}
	for e := range Merge(ctx, []*Source{a, b}) {
		names[e.Source] = true
	}
	if !names["src-a"] || !names["src-b"] {
		t.Errorf("source names not preserved: %v", names)
	}
}

func TestMerge_Empty(t *testing.T) {
	ctx := context.Background()
	var got []Entry
	for e := range Merge(ctx, nil) {
		got = append(got, e)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestMerge_LinesContent(t *testing.T) {
	a := New("a", nopCloser("x\ny"))
	ctx := context.Background()
	var lines []string
	for e := range Merge(ctx, []*Source{a}) {
		lines = append(lines, e.Line)
	}
	sort.Strings(lines)
	if lines[0] != "x" || lines[1] != "y" {
		t.Errorf("unexpected lines: %v", lines)
	}
}
