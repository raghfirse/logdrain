package formatter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBookmark_GetDefault(t *testing.T) {
	dir := t.TempDir()
	b, err := NewBookmark(filepath.Join(dir, "bm.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := b.Get("app"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestBookmark_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	b, err := NewBookmark(filepath.Join(dir, "bm.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Set("app", 42); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if got := b.Get("app"); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestBookmark_PersistsAcrossReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bm.json")

	b1, _ := NewBookmark(path)
	_ = b1.Set("svc", 99)

	b2, err := NewBookmark(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := b2.Get("svc"); got != 99 {
		t.Errorf("expected 99, got %d", got)
	}
}

func TestBookmark_Reset(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bm.json")

	b, _ := NewBookmark(path)
	_ = b.Set("x", 7)
	if err := b.Reset(); err != nil {
		t.Fatalf("Reset error: %v", err)
	}
	if got := b.Get("x"); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed")
	}
}

func TestBookmark_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bm.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := NewBookmark(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseBookmarkFlag_Valid(t *testing.T) {
	f, err := ParseBookmarkFlag("/tmp/logdrain.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Path != "/tmp/logdrain.json" {
		t.Errorf("unexpected path: %s", f.Path)
	}
}

func TestParseBookmarkFlag_Empty(t *testing.T) {
	_, err := ParseBookmarkFlag("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestBookmarkFlag_String_Default(t *testing.T) {
	var f BookmarkFlag
	if f.String() != "(none)" {
		t.Errorf("expected (none), got %s", f.String())
	}
}

func TestBookmarkFlag_Type(t *testing.T) {
	var f BookmarkFlag
	if f.Type() != "bookmark" {
		t.Errorf("expected bookmark, got %s", f.Type())
	}
}
