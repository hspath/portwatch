package ports

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRateLimiter_FirstCallAllowed(t *testing.T) {
	rl := NewRateLimiter(5 * time.Second)
	if !rl.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRateLimiter_SecondCallSuppressed(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(5 * time.Second)
	rl.nowFunc = fixedClock(base)

	rl.Allow("tcp:8080")
	if rl.Allow("tcp:8080") {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestRateLimiter_AllowedAfterWindowExpires(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(5 * time.Second)
	rl.nowFunc = fixedClock(base)
	rl.Allow("tcp:8080")

	// advance past the window
	rl.nowFunc = fixedClock(base.Add(6 * time.Second))
	if !rl.Allow("tcp:8080") {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestRateLimiter_DifferentKeysIndependent(t *testing.T) {
	rl := NewRateLimiter(5 * time.Second)
	rl.Allow("tcp:8080")

	if !rl.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestRateLimiter_Purge_RemovesExpiredEntries(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(5 * time.Second)
	rl.nowFunc = fixedClock(base)

	rl.Allow("tcp:8080")
	rl.Allow("tcp:9090")

	if rl.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", rl.Len())
	}

	// advance past window and purge
	rl.nowFunc = fixedClock(base.Add(6 * time.Second))
	rl.Purge()

	if rl.Len() != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", rl.Len())
	}
}

func TestRateLimiter_Purge_KeepsActiveEntries(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(5 * time.Second)

	rl.nowFunc = fixedClock(base)
	rl.Allow("tcp:8080")

	// second key added later, still within window at purge time
	rl.nowFunc = fixedClock(base.Add(4 * time.Second))
	rl.Allow("tcp:9090")

	// purge at t+6: tcp:8080 expired, tcp:9090 still active
	rl.nowFunc = fixedClock(base.Add(6 * time.Second))
	rl.Purge()

	if rl.Len() != 1 {
		t.Fatalf("expected 1 active entry after purge, got %d", rl.Len())
	}
}
