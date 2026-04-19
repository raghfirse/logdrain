package formatter

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RateLimit holds parsed rate limit configuration.
type RateLimit struct {
	Max    int
	Window time.Duration
}

// ParseRateFlag parses a rate limit string of the form "N/duration", e.g. "10/1s".
func ParseRateFlag(s string) (RateLimit, error) {
	if s == "" {
		return RateLimit{}, fmt.Errorf("rate flag must not be empty")
	}
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return RateLimit{}, fmt.Errorf("rate flag must be in format N/duration, got %q", s)
	}
	n, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || n < 0 {
		return RateLimit{}, fmt.Errorf("invalid rate count %q", parts[0])
	}
	d, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil || d <= 0 {
		return RateLimit{}, fmt.Errorf("invalid rate window %q", parts[1])
	}
	return RateLimit{Max: n, Window: d}, nil
}

// rateLimitFlag implements flag.Value for --rate.
type rateLimitFlag struct {
	value *RateLimit
}

func (f *rateLimitFlag) Set(s string) error {
	rl, err := ParseRateFlag(s)
	if err != nil {
		return err
	}
	*f.value = rl
	return nil
}

func (f *rateLimitFlag) String() string {
	if f.value == nil || f.value.Window == 0 {
		return ""
	}
	return fmt.Sprintf("%d/%s", f.value.Max, f.value.Window)
}

func (f *rateLimitFlag) Type() string { return "rate" }
