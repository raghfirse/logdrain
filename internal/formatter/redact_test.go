package formatter

import (
	"encoding/json"
	"testing"
)

func TestRedactor_NoKeys(t *testing.T) {
	r := NewRedactor(nil)
	line := `{"password":"secret","user":"alice"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestRedactor_MasksField(t *testing.T) {
	r := NewRedactor([]string{"password"})
	line := `{"password":"secret","user":"alice"}`
	out := r.Apply(line)
	var obj map[string]string
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if obj["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %s", obj["password"])
	}
	if obj["user"] != "alice" {
		t.Errorf("user should be unchanged")
	}
}

func TestRedactor_CaseInsensitive(t *testing.T) {
	r := NewRedactor([]string{"TOKEN"})
	line := `{"token":"abc123"}`
	out := r.Apply(line)
	var obj map[string]string
	json.Unmarshal([]byte(out), &obj)
	if obj["token"] != "[REDACTED]" {
		t.Errorf("expected redacted, got %s", obj["token"])
	}
}

func TestRedactor_NonJSONUnchanged(t *testing.T) {
	r := NewRedactor([]string{"password"})
	line := "not json at all"
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged non-json line")
	}
}

func TestRedactor_NoMatchingFields(t *testing.T) {
	r := NewRedactor([]string{"secret"})
	line := `{"user":"alice","level":"info"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged when no fields match")
	}
}

func TestParseRedactFlag_Valid(t *testing.T) {
	keys := ParseRedactFlag("password, token, secret")
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[1] != "token" {
		t.Errorf("expected token, got %s", keys[1])
	}
}

func TestParseRedactFlag_Empty(t *testing.T) {
	if keys := ParseRedactFlag(""); keys != nil {
		t.Errorf("expected nil for empty input")
	}
}
