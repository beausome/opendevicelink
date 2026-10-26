package lease

import (
	"errors"
	"testing"
	"time"
)

func TestLeaseIsExclusiveAndReusableAfterRelease(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	manager := NewManager(func() time.Time { return now })
	first, err := manager.Acquire("camera-1", "alice", time.Minute)
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	if _, err := manager.Acquire("camera-1", "bob", time.Minute); !errors.Is(err, ErrBusy) {
		t.Fatalf("second acquire error = %v, want ErrBusy", err)
	}
	if err := manager.Release(first.Token); err != nil {
		t.Fatalf("release: %v", err)
	}
	if _, err := manager.Acquire("camera-1", "bob", time.Minute); err != nil {
		t.Fatalf("acquire after release: %v", err)
	}
}

func TestExpiredLeaseIsReclaimed(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	manager := NewManager(func() time.Time { return now })
	if _, err := manager.Acquire("printer-1", "alice", time.Second); err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	now = now.Add(2 * time.Second)
	if _, err := manager.Acquire("printer-1", "bob", time.Minute); err != nil {
		t.Fatalf("expired lease was not reclaimed: %v", err)
	}
}
