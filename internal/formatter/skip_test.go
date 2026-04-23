package formatter

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSkip_ZeroSkipsNothing(t *testing.T) {
	sk, err := NewSkip(0, "")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		line := fmt.Sprintf(`{"n":%d}`, i)
		out, ok := sk.Apply(line)
		if !ok || out != line {
			t.Errorf("line %d: expected pass-through", i)
		}
	}
}

func TestSkip_SkipsFirstN(t *testing.T) {
	sk, _ := NewSkip(3, "")
	lines := []string{"a", "b", "c", "d", "e"}
	expected := []bool{false, false, false, true, true}
	for i, l := range lines {
		_, ok := sk.Apply(l)
		if ok != expected[i] {
			t.Errorf("line %d: got ok=%v, want %v", i, ok, expected[i])
		}
	}
}

func TestSkip_PerKeyField(t *testing.T) {
	sk, _ := NewSkip(2, "svc")
	makeLine := func(svc string, n int) string {
		b, _ := json.Marshal(map[string]interface{}{"svc": svc, "n": n})
		return string(b)
	}
	// First two lines for "a" are skipped.
	if _, ok := sk.Apply(makeLine("a", 1)); ok {
		t.Error("a/1: should be skipped")
	}
	if _, ok := sk.Apply(makeLine("a", 2)); ok {
		t.Error("a/2: should be skipped")
	}
	// Third line for "a" passes.
	if _, ok := sk.Apply(makeLine("a", 3)); !ok {
		t.Error("a/3: should pass")
	}
	// First two lines for "b" are skipped independently.
	if _, ok := sk.Apply(makeLine("b", 1)); ok {
		t.Error("b/1: should be skipped")
	}
	if _, ok := sk.Apply(makeLine("b", 2)); ok {
		t.Error("b/2: should be skipped")
	}
	if _, ok := sk.Apply(makeLine("b", 3)); !ok {
		t.Error("b/3: should pass")
	}
}

func TestSkip_Reset(t *testing.T) {
	sk, _ := NewSkip(2, "")
	sk.Apply("line1")
	sk.Apply("line2")
	_, ok := sk.Apply("line3")
	if !ok {
		t.Fatal("line3 should pass before reset")
	}
	sk.Reset()
	if _, ok := sk.Apply("line1"); ok {
		t.Error("after reset, first line should be skipped again")
	}
}

func TestSkip_NegativeLimitError(t *testing.T) {
	_, err := NewSkip(-1, "")
	if err == nil {
		t.Error("expected error for negative limit")
	}
}

func TestParseSkipFlag_Valid(t *testing.T) {
	cases := []struct {
		input    string
		wantN    int
		wantKey  string
	}{
		{"5", 5, ""},
		{"0", 0, ""},
		{"10:service", 10, "service"},
		{"", 0, ""},
	}
	for _, c := range cases {
		sk, err := ParseSkipFlag(c.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.input, err)
			continue
		}
		if sk.limit != c.wantN || sk.keyField != c.wantKey {
			t.Errorf("%q: got limit=%d key=%q, want limit=%d key=%q",
				c.input, sk.limit, sk.keyField, c.wantN, c.wantKey)
		}
	}
}

func TestParseSkipFlag_Invalid(t *testing.T) {
	for _, bad := range []string{"abc", "x:field", "-3"} {
		if _, err := ParseSkipFlag(bad); err == nil {
			t.Errorf("%q: expected error", bad)
		}
	}
}
