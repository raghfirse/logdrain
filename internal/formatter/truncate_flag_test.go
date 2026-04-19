package formatter

import "testing"

func TestTruncateFlag_Set_Valid(t *testing.T) {
	var f TruncateFlag
	if err := f.Set("100"); err != nil {
		t.Fatal(err)
	}
	if f.MaxLen != 100 {
		t.Fatalf("expected 100, got %d", f.MaxLen)
	}
}

func TestTruncateFlag_Set_Invalid(t *testing.T) {
	var f TruncateFlag
	if err := f.Set("bad"); err == nil {
		t.Fatal("expected error")
	}
}

func TestTruncateFlag_String_Default(t *testing.T) {
	var f TruncateFlag
	if f.String() != "0" {
		t.Fatalf("expected 0, got %s", f.String())
	}
}

func TestTruncateFlag_String_AfterSet(t *testing.T) {
	var f TruncateFlag
	_ = f.Set("50")
	if f.String() != "50:…" {
		t.Fatalf("unexpected: %s", f.String())
	}
}

func TestTruncateFlag_Type(t *testing.T) {
	var f TruncateFlag
	if f.Type() != "int[:suffix]" {
		t.Fatalf("unexpected type: %s", f.Type())
	}
}

func TestTruncateFlag_Truncator_Nil(t *testing.T) {
	var f TruncateFlag
	if f.Truncator() != nil {
		t.Fatal("expected nil truncator when disabled")
	}
}

func TestTruncateFlag_Truncator_Active(t *testing.T) {
	var f TruncateFlag
	_ = f.Set("10")
	tr := f.Truncator()
	if tr == nil {
		t.Fatal("expected non-nil truncator")
	}
	out := tr.Apply(map[string]any{"msg": "hello world!"})
	if out["msg"] != "hello worl…" {
		t.Fatalf("unexpected: %v", out["msg"])
	}
}
