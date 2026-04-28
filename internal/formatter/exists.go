package formatter

import (
	"encoding/json"
	"strings"
)

// Exists filters log lines based on whether specified fields are present (or absent).
type Exists struct {
	required []string
	forbidden []string
}

// ExistsConfig holds the configuration for an Exists filter.
type ExistsConfig struct {
	// Required fields that must be present in the log line.
	Required []string
	// Forbidden fields that must not be present in the log line.
	Forbidden []string
}

// NewExists creates a new Exists filter.
// Returns an error if no fields are specified.
func NewExists(cfg ExistsConfig) (*Exists, error) {
	if len(cfg.Required) == 0 && len(cfg.Forbidden) == 0 {
		return &Exists{}, nil
	}
	return &Exists{
		required:  cfg.Required,
		forbidden: cfg.Forbidden,
	}, nil
}

// Apply returns the line if it satisfies the field existence rules,
// or an empty string if it should be suppressed.
func (e *Exists) Apply(line string) string {
	if len(e.required) == 0 && len(e.forbidden) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, key := range e.required {
		if !hasKeyCaseInsensitive(obj, key) {
			return ""
		}
	}

	for _, key := range e.forbidden {
		if hasKeyCaseInsensitive(obj, key) {
			return ""
		}
	}

	return line
}

func hasKeyCaseInsensitive(obj map[string]json.RawMessage, key string) bool {
	lower := strings.ToLower(key)
	for k := range obj {
		if strings.ToLower(k) == lower {
			return true
		}
	}
	return false
}

// ParseExistsFlag parses a raw flag string into an ExistsConfig.
// Format: "field1,field2" for required; prefix with "!" for forbidden.
// Example: "user,!debug" — requires "user", forbids "debug".
func ParseExistsFlag(raw string) (ExistsConfig, error) {
	var cfg ExistsConfig
	if raw == "" {
		return cfg, nil
	}
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "!") {
			cfg.Forbidden = append(cfg.Forbidden, p[1:])
		} else {
			cfg.Required = append(cfg.Required, p)
		}
	}
	return cfg, nil
}
