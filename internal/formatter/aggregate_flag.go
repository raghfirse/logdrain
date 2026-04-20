package formatter

import (
	"fmt"
	"strings"
)

// AggregateFlagValue implements flag.Value for --aggregate field:mode.
type AggregateFlagValue struct {
	Field string
	Mode  AggregateMode
	set   bool
}

func (f *AggregateFlagValue) String() string {
	if !f.set {
		return ""
	}
	return fmt.Sprintf("%s:%s", f.Field, aggregateModeName(f.Mode))
}

func (f *AggregateFlagValue) Set(s string) error {
	field, mode, err := ParseAggregateFlag(s)
	if err != nil {
		return err
	}
	f.Field = field
	f.Mode = mode
	f.set = true
	return nil
}

func (f *AggregateFlagValue) Type() string { return "field:mode" }

// ParseAggregateFlag parses "field:mode" into field name and AggregateMode.
func ParseAggregateFlag(s string) (string, AggregateMode, error) {
	if s == "" {
		return "", 0, fmt.Errorf("aggregate flag must not be empty")
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("aggregate flag must be field:mode, got %q", s)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return "", 0, fmt.Errorf("aggregate field must not be empty")
	}
	mode, err := parseAggregateMode(strings.TrimSpace(parts[1]))
	if err != nil {
		return "", 0, err
	}
	return field, mode, nil
}

func parseAggregateMode(s string) (AggregateMode, error) {
	switch strings.ToLower(s) {
	case "count":
		return AggregateCount, nil
	case "sum":
		return AggregateSum, nil
	case "min":
		return AggregateMin, nil
	case "max":
		return AggregateMax, nil
	default:
		return 0, fmt.Errorf("unknown aggregate mode %q: must be count, sum, min, or max", s)
	}
}

func aggregateModeName(m AggregateMode) string {
	switch m {
	case AggregateCount:
		return "count"
	case AggregateSum:
		return "sum"
	case AggregateMin:
		return "min"
	case AggregateMax:
		return "max"
	default:
		return "unknown"
	}
}
