package formatter

import "testing"

func TestSampler_Disabled(t *testing.T) {
	s := NewSampler(0)
	for i := 0; i < 10; i++ {
		if !s.Sample("k") {
			t.Fatal("expected all lines to pass when n=0")
		}
	}
}

func TestSampler_EveryN(t *testing.T) {
	s := NewSampler(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Sample("k")
	}
	expected := []bool{true, false, false, true, false, false, true, false, false}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: got %v want %v", i, got, expected[i])
		}
	}
}

func TestSampler_SeparateKeys(t *testing.T) {
	s := NewSampler(2)
	if !s.Sample("a") {
		t.Error("first call for 'a' should pass")
	}
	if !s.Sample("b") {
		t.Error("first call for 'b' should pass")
	}
	if s.Sample("a") {
		t.Error("second call for 'a' should not pass")
	}
}

func TestSampler_Reset(t *testing.T) {
	s := NewSampler(2)
	s.Sample("k") // 1 -> pass
	s.Sample("k") // 2 -> drop
	s.Reset()
	if !s.Sample("k") {
		t.Error("after reset first call should pass")
	}
}
