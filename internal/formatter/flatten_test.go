package formatter

import (
	"testing"
)

func TestFlattener_NonJSON(t *testing.T) {
	f := NewFlattener(".")
	input := "not json"
	if got := f.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestFlattener_AlreadyFlat(t *testing.T) {
	f := NewFlattener(".")
	input := `{"a":1,"b":"hello"}`
	got := f.Apply(input)
	assertJSONField(t, got, "a", float64(1))
	assertJSONField(t, got, "b", "hello")
}

func TestFlattener_NestedObject(t *testing.T) {
	f := NewFlattener(".")
	input := `{"a":{"b":{"c":42}}}`
	got := f.Apply(input)
	assertJSONField(t, got, "a.b.c", float64(42))
}

func TestFlattener_CustomSeparator(t *testing.T) {
	f := NewFlattener("_")
	input := `{"x":{"y":"val"}}`
	got := f.Apply(input)
	assertJSONField(t, got, "x_y", "val")
}

func TestFlattener_DefaultSeparator(t *testing.T) {
	f := NewFlattener("")
	if f.separator != "." {
		t.Errorf("expected default separator '.', got %q", f.separator)
	}
}

func TestFlattener_MixedDepth(t *testing.T) {
	f := NewFlattener(".")
	input := `{"top":1,"nested":{"mid":2,"deep":{"val":3}}}`
	got := f.Apply(input)
	assertJSONField(t, got, "top", float64(1))
	assertJSONField(t, got, "nested.mid", float64(2))
	assertJSONField(t, got, "nested.deep.val", float64(3))
}

func TestParseFlattenFlag_Valid(t *testing.T) {
	f, err := ParseFlattenFlag(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.separator != "." {
		t.Errorf("expected '.', got %q", f.separator)
	}
}

func TestParseFlattenFlag_Empty(t *testing.T) {
	f, err := ParseFlattenFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.separator != "." {
		t.Errorf("expected default '.', got %q", f.separator)
	}
}

func TestParseFlattenFlag_WhitespaceSeparator(t *testing.T) {
	_, err := ParseFlattenFlag(" ")
	if err == nil {
		t.Error("expected error for whitespace separator")
	}
}
