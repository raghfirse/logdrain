package formatter

import (
	"encoding/json"
	"testing"
)

func TestJoin_NoRules(t *testing.T) {
	j := NewJoin(nil)
	line := `{"a":"hello","b":"world"}`
	if got := j.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestJoin_CombinesTwoFields(t *testing.T) {
	j := NewJoin([]JoinRule{{Keys: []string{"first", "last"}, OutputKey: "full", Sep: " "}})
	line := `{"first":"John","last":"Doe"}`
	out := j.Apply(line)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["full"] != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", obj["full"])
	}
}

func TestJoin_CustomSeparator(t *testing.T) {
	j := NewJoin([]JoinRule{{Keys: []string{"host", "port"}, OutputKey: "addr", Sep: ":"}})
	line := `{"host":"localhost","port":"8080"}`
	out := j.Apply(line)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["addr"] != "localhost:8080" {
		t.Errorf("expected 'localhost:8080', got %q", obj["addr"])
	}
}

func TestJoin_MissingKeySkipped(t *testing.T) {
	j := NewJoin([]JoinRule{{Keys: []string{"a", "b", "c"}, OutputKey: "out", Sep: "-"}})
	line := `{"a":"x","c":"z"}`
	out := j.Apply(line)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["out"] != "x-z" {
		t.Errorf("expected 'x-z', got %q", obj["out"])
	}
}

func TestJoin_NonJSONPassthrough(t *testing.T) {
	j := NewJoin([]JoinRule{{Keys: []string{"a", "b"}, OutputKey: "c", Sep: " "}})
	line := "not json"
	if got := j.Apply(line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestParseJoinFlag_Valid(t *testing.T) {
	rule, err := ParseJoinFlag("first+last-> full_name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.OutputKey != "full_name" {
		t.Errorf("expected output 'full_name', got %q", rule.OutputKey)
	}
	if len(rule.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(rule.Keys))
	}
	if rule.Sep != " " {
		t.Errorf("expected sep ' ', got %q", rule.Sep)
	}
}

func TestParseJoinFlag_WithSep(t *testing.T) {
	rule, err := ParseJoinFlag("host+port:separator->addr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Sep != "separator" {
		t.Errorf("expected sep 'separator', got %q", rule.Sep)
	}
}

func TestParseJoinFlag_Invalid(t *testing.T) {
	cases := []string{
		"no_arrow",
		"only_one->out",
		"+b->out",
		"a+b->",
	}
	for _, c := range cases {
		if _, err := ParseJoinFlag(c); err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}
