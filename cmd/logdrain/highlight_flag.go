package main

import (
	"fmt"
	"strings"
)

// highlightFlag implements flag.Value for repeated --highlight field=color flags.
type highlightFlag []string

func (h *highlightFlag) String() string {
	if h == nil {
		return ""
	}
	return strings.Join(*h, ", ")
}

func (h *highlightFlag) Set(val string) error {
	if !strings.Contains(val, "=") {
		return fmt.Errorf("highlight flag %q must be in field=color format", val)
	}
	*h = append(*h, val)
	return nil
}
