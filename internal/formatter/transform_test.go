package formatter

import (
	"testing"
)

func TestNewTransform_Valid(t *testing.T) {
	_, err := NewTransform([]string{"msg:message", "ts:timestamp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewTransform_Invalid(t *testing.T) {
	cases := []string{"nocodon", ":empty", "empty:", ""}
	for _, c := range cases {
		_, err := NewTransform([]string{c})
		if err == nil {
			t.Errorf("expected error for rule %q", c)
		}
	}
}

func TestTransform_Apply_RenamesField(t *testing.T) {
	tr, _ := NewTransform([]string{"msg:message"})
	out := tr.Apply(`{"msg":"hello","level":"info"}`)
	if !containsKey(out, "message") {
		t.Errorf("expected key 'message' in output: %s", out)
	}
	if containsKey(out, "msg") {
		t.Errorf("expected key 'msg' to be removed: %s", out)
	}
}

func TestTransform_Apply_NoRules(t *testing.T) {
	tr, _ := NewTransform(nil)
	line := `{"msg":"hello"}`
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %s", got)
	}
}

func TestTransform_Apply_NonJSON(t *testing.T) {
	tr, _ := NewTransform([]string{"msg:message"})
	line := "plain text line"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged non-JSON line, got %s", got)
	}
}

func TestTransform_Apply_CaseInsensitive(t *testing.T) {
	tr, _ := NewTransform([]string{"MSG:message"})
	out := tr.Apply(`{"msg":"hello"}`)
	if !containsKey(out, "message") {
		t.Errorf("expected key 'message' in output: %s", out)
	}
}

func TestTransform_Apply_MissingKey(t *testing.T) {
	tr, _ := NewTransform([]string{"missing:target"})
	line := `{"level":"info"}`
	out := tr.Apply(line)
	if containsKey(out, "target") {
		t.Errorf("expected no 'target' key when source missing: %s", out)
	}
}

func TestTransform_Apply_MultipleRules(t *testing.T) {
	tr, _ := NewTransform([]string{"msg:message", "ts:timestamp"})
	out := tr.Apply(`{"msg":"hi","ts":"2024-01-01T00:00:00Z"}`)
	if !containsKey(out, "message") || !containsKey(out, "timestamp") {
		t.Errorf("expected both renamed keys in output: %s", out)
	}
}

func containsKey(json, key string) bool {
	return len(json) > 0 && (len(key) > 0) &&
		(func() bool {
			for i := 0; i < len(json)-len(key)-1; i++ {
				if json[i] == '"' && json[i+1:i+1+len(key)] == key && json[i+1+len(key)] == '"' {
					return true
				}
			}
			return false
		})()
}
