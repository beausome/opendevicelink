package simulator

import (
	"errors"
	"testing"
)

func TestAttachLifecycle(t *testing.T) {
	backend := New()
	if err := backend.Attach("sim-camera-1"); err != nil {
		t.Fatalf("attach: %v", err)
	}
	if err := backend.Attach("sim-camera-1"); err == nil {
		t.Fatal("duplicate attach accepted")
	}
	if err := backend.Detach("sim-camera-1"); err != nil {
		t.Fatalf("detach: %v", err)
	}
}

func TestUnknownDeviceIsRejected(t *testing.T) {
	if err := New().Attach("missing"); !errors.Is(err, ErrDeviceNotFound) {
		t.Fatalf("error = %v, want ErrDeviceNotFound", err)
	}
}
