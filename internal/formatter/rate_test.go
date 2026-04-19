package formatter

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowsUnderLimit(t *testing.T) {
	r := NewRateLimiter(3, time.Second)
	for i := 0; i < 3; i++ {
		if !r.Allow("key") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	r := NewRateLimiter(2, time.Second)
	r.Allow("key")
	r.Allow("key")
	if r.Allow("key") {
		t.Fatal("expected Allow to return false when over limit")
	}
}

func TestRateLimiter_SeparateKeys(t *testing.T) {
	r := NewRateLimiter(1, time.Second)
	if !r.Allow("a") {
		t.Fatal("expected true for key a")
	}
	if !r.Allow("b") {
		t.Fatal("expected true for key b")
	}
	if r.Allow("a") {
		t.Fatal("expected false for key a after limit")
	}
}

func TestRateLimiter_ExpiresAfterWindow(t *testing.T) {
	now := time.Now()
	r := NewRateLimiter(1, time.Second)
	r.nowFn = func() time.Time { return now }
	r.Allow("key")

	r.nowFn = func() time.Time { return now.Add(2 * time.Second) }
	if !r.Allow("key") {
		t.Fatal("expected Allow to return true after window expires")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	r := NewRateLimiter(1, time.Second)
	r.Allow("key")
	r.Reset()
	if !r.Allow("key") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestRateLimiter_ZeroMax(t *testing.T) {
	r := NewRateLimiter(0, time.Second)
	if r.Allow("key") {
		t.Fatal("expected Allow to return false with max=0")
	}
}
