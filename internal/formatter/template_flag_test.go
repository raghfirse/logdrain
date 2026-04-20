package formatter

import "testing"

func TestTemplateFlag_Set_Valid(t *testing.T) {
	var f TemplateFlag
	err := f.Set(`{{index . "level"}}: {{index . "msg"}}`)
	if err != nil {
		t.Fatalf("Set() unexpected error: %v", err)
	}
	if !f.IsSet() {
		t.Error("expected IsSet() to be true after Set()")
	}
	if f.Template() == nil {
		t.Error("expected non-nil Template() after Set()")
	}
}

func TestTemplateFlag_Set_Invalid(t *testing.T) {
	var f TemplateFlag
	err := f.Set("{{unclosed")
	if err == nil {
		t.Fatal("Set() expected error for invalid template")
	}
	if f.IsSet() {
		t.Error("expected IsSet() to be false after failed Set()")
	}
}

func TestTemplateFlag_String_Default(t *testing.T) {
	var f TemplateFlag
	if got := f.String(); got != "" {
		t.Errorf("String() = %q, want empty string", got)
	}
}

func TestTemplateFlag_String_AfterSet(t *testing.T) {
	var f TemplateFlag
	input := `{{index . "msg"}}`
	_ = f.Set(input)
	if got := f.String(); got != input {
		t.Errorf("String() = %q, want %q", got, input)
	}
}

func TestTemplateFlag_Type(t *testing.T) {
	var f TemplateFlag
	if got := f.Type(); got != "template" {
		t.Errorf("Type() = %q, want \"template\"", got)
	}
}

func TestParseTemplateFlagValue_Valid(t *testing.T) {
	_, err := ParseTemplateFlagValue(`{{index . "level"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseTemplateFlagValue_Invalid(t *testing.T) {
	_, err := ParseTemplateFlagValue("{{bad")
	if err == nil {
		t.Fatal("expected error")
	}
}
