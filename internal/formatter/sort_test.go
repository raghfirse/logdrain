package formatter

import (
	"testing"
)

func TestNewSorter_EmptyField(t *testing.T) {
	_, err := NewSorter("", "asc")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSorter_InvalidDirection(t *testing.T) {
	_, err := NewSorter("level", "sideways")
	if err == nil {
		t.Fatal("expected error for invalid direction")
	}
}

func TestSorter_AscendingOrder(t *testing.T) {
	s, _ := NewSorter("ts", "asc")
	lines := []string{
		`{"ts":"2024-01-03","msg":"c"}`,
		`{"ts":"2024-01-01","msg":"a"}`,
		`{"ts":"2024-01-02","msg":"b"}`,
	}
	for _, l := range lines {
		if out := s.Apply(l); out != nil {
			t.Fatalf("Apply should return nil, got %v", out)
		}
	}
	out := s.Flush()
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
	for i, want := range []string{"a", "b", "c"} {
		var obj map[string]string
		if err := unmarshalJSON(out[i], &obj); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if obj["msg"] != want {
			t.Errorf("line %d: got msg=%q, want %q", i, obj["msg"], want)
		}
	}
}

func TestSorter_DescendingOrder(t *testing.T) {
	s, _ := NewSorter("ts", "desc")
	lines := []string{
		`{"ts":"2024-01-01","msg":"a"}`,
		`{"ts":"2024-01-03","msg":"c"}`,
		`{"ts":"2024-01-02","msg":"b"}`,
	}
	for _, l := range lines {
		s.Apply(l)
	}
	out := s.Flush()
	for i, want := range []string{"c", "b", "a"} {
		var obj map[string]string
		if err := unmarshalJSON(out[i], &obj); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if obj["msg"] != want {
			t.Errorf("line %d: got msg=%q, want %q", i, obj["msg"], want)
		}
	}
}

func TestSorter_FlushClearsBuffer(t *testing.T) {
	s, _ := NewSorter("level", "asc")
	s.Apply(`{"level":"info"}`)
	s.Flush()
	out := s.Flush()
	if len(out) != 0 {
		t.Fatalf("expected empty flush after reset, got %d lines", len(out))
	}
}

func TestSorter_NonJSONBuffered(t *testing.T) {
	s, _ := NewSorter("level", "asc")
	s.Apply("not json")
	s.Apply(`{"level":"info"}`)
	out := s.Flush()
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
}

func TestParseSortFlag_Valid(t *testing.T) {
	s, err := ParseSortFlag("ts:desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.descending {
		t.Error("expected descending")
	}
}

func TestParseSortFlag_DefaultAsc(t *testing.T) {
	s, err := ParseSortFlag("ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.descending {
		t.Error("expected ascending by default")
	}
}

func TestParseSortFlag_Empty(t *testing.T) {
	_, err := ParseSortFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}

// helper reused across tests
func unmarshalJSON(s string, v interface{}) error {
	import_json_bytes := []byte(s)
	_ = import_json_bytes
	return jsonUnmarshal([]byte(s), v)
}
