package formatter

import (
	"testing"
	"time"
)

func TestSinceFlag_Set_Valid(t *testing.T) {
	var f SinceFlag
	if err := f.Set("10m"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.raw != "10m" {
		t.Errorf("expected raw=10m, got %q", f.raw)
	}
	if f.Cutoff.IsZero() {
		t.Error("expected non-zero cutoff")
	}
	expected := time.Now().Add(-10 * time.Minute)
	if f.Cutoff.After(expected.Add(time.Second)) || f.Cutoff.Before(expected.Add(-time.Second)) {
		t.Errorf("cutoff %v not near expected %v", f.Cutoff, expected)
	}
}

func TestSinceFlag_Set_Invalid(t *testing.T) {
	var f SinceFlag
	if err := f.Set("badvalue"); err == nil {
		t.Error("expected error for invalid duration")
	}
}

func TestSinceFlag_Set_Negative(t *testing.T) {
	var f SinceFlag
	if err := f.Set("-1h"); err == nil {
		t.Error("expected error for negative duration")
	}
}

func TestSinceFlag_String_Default(t *testing.T) {
	var f SinceFlag
	if s := f.String(); s != "" {
		t.Errorf("expected empty string default, got %q", s)
	}
}

func TestSinceFlag_String_AfterSet(t *testing.T) {
	var f SinceFlag
	_ = f.Set("2h")
	if s := f.String(); s != "2h" {
		t.Errorf("expected \"2h\", got %q", s)
	}
}

func TestSinceFlag_Type(t *testing.T) {
	var f SinceFlag
	if f.Type() != "duration" {
		t.Errorf("expected type=duration, got %q", f.Type())
	}
}

func TestSinceFlag_IsSet(t *testing.T) {
	var f SinceFlag
	if f.IsSet() {
		t.Error("expected IsSet=false before Set")
	}
	_ = f.Set("30s")
	if !f.IsSet() {
		t.Error("expected IsSet=true after Set")
	}
}
