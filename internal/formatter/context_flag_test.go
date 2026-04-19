package formatter

import "testing"

func TestContextFlag_Set_Valid(t *testing.T) {
	var f ContextFlag
	if err := f.Set("request_id,user"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Fields()) != 2 {
		t.Errorf("expected 2 fields, got %d", len(f.Fields()))
	}
}

func TestContextFlag_Set_Empty(t *testing.T) {
	var f ContextFlag
	if err := f.Set(""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Fields()) != 0 {
		t.Errorf("expected 0 fields, got %d", len(f.Fields()))
	}
}

func TestContextFlag_String_Default(t *testing.T) {
	var f ContextFlag
	if s := f.String(); s != "" {
		t.Errorf("expected empty string, got %q", s)
	}
}

func TestContextFlag_String_AfterSet(t *testing.T) {
	var f ContextFlag
	_ = f.Set("a,b,c")
	if s := f.String(); s != "a,b,c" {
		t.Errorf("expected 'a,b,c', got %q", s)
	}
}

func TestContextFlag_Type(t *testing.T) {
	var f ContextFlag
	if f.Type() != "fields" {
		t.Errorf("expected 'fields', got %q", f.Type())
	}
}
