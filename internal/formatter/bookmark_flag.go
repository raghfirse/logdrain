package formatter

import (
	"fmt"
	"strings"
)

// BookmarkFlag is a flag value that holds the path to a bookmark file.
type BookmarkFlag struct {
	Path string
}

// Set parses and validates the bookmark file path.
func (f *BookmarkFlag) Set(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("bookmark path must not be empty")
	}
	f.Path = s
	return nil
}

// String returns the current bookmark path or a default label.
func (f *BookmarkFlag) String() string {
	if f.Path == "" {
		return "(none)"
	}
	return f.Path
}

// Type returns the flag type name for use with pflag/flag.
func (f *BookmarkFlag) Type() string {
	return "bookmark"
}

// ParseBookmarkFlag parses a raw string into a BookmarkFlag.
func ParseBookmarkFlag(s string) (BookmarkFlag, error) {
	var f BookmarkFlag
	if err := f.Set(s); err != nil {
		return BookmarkFlag{}, err
	}
	return f, nil
}
