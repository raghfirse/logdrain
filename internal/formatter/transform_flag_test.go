package formatter

import (
	"testing"
)

func TestTransformFlag_Set_Valid(t *testing.T) {
	var f TransformFlag
	if err := f.Set("msg:message"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules()) != 1 || f.Rules()[0] != "msg:message" {
		t.Errorf("unexpected rules: %v", f.Rules())
	}
}

func TestTransformFlag_Set_Invalid(t *testing.T) {
	cases := []string{"nodivider", ":empty", "empty:", ""}
	for _, c := range cases {
		var f TransformFlag
		if err := f.Set(c); err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestTransformFlag_String_Default(t *testing.T) {
	var f TransformFlag
	if f.String() != "none" {
		t.Errorf("expected 'none', got %q", f.String())
	}
}

func TestTransformFlag_String_AfterSet(t *testing.T) {
	var f TransformFlag
	_ = f.Set("msg:message")
	_ = f.Set("ts:timestamp")
	got := f.String()
	if got != "msg:message,ts:timestamp" {
		t.Errorf("unexpected string: %q", got)
	}
}

func TestTransformFlag_Type(t *testing.T) {
	var f TransformFlag
	if f.Type() != "from:to" {
		t.Errorf("unexpected type: %q", f.Type())
	}
}

func TestParseTransformFlag_Nil(t *testing.T) {
	tr, err := ParseTransformFlag(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Transform")
	}
}

func TestParseTransformFlag_WithRules(t *testing.T) {
	var f TransformFlag
	_ = f.Set("msg:message")
	tr, err := ParseTransformFlag(&f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tr.Apply(`{"msg":"hello"}`)
	if !containsKey(out, "message") {
		t.Errorf("expected 'message' key in output: %s", out)
	}
}
