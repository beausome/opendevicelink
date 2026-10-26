// Package simulator provides a deterministic backend for protocol development.
package simulator

import (
	"errors"
	"sync"

	"github.com/beausome/opendevicelink/internal/device"
)

var ErrDeviceNotFound = errors.New("device not found")

type Backend struct {
	mu       sync.Mutex
	devices  []device.Descriptor
	attached map[string]bool
}

func New() *Backend {
	return &Backend{
		devices: []device.Descriptor{
			{ID: "sim-camera-1", Name: "Conference Camera", VendorID: 0x1209, ProductID: 0x0001, Class: device.ClassCamera},
			{ID: "sim-keyboard-1", Name: "Mechanical Keyboard", VendorID: 0x1209, ProductID: 0x0002, Class: device.ClassHID},
			{ID: "sim-drive-1", Name: "Backup Drive", VendorID: 0x1209, ProductID: 0x0003, Class: device.ClassMassStorage},
		},
		attached: make(map[string]bool),
	}
}

func DefaultPolicy() device.Policy {
	return device.DefaultPolicy()
}

func (b *Backend) List() ([]device.Descriptor, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	result := make([]device.Descriptor, len(b.devices))
	copy(result, b.devices)
	return result, nil
}

func (b *Backend) Attach(deviceID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.exists(deviceID) {
		return ErrDeviceNotFound
	}
	if b.attached[deviceID] {
		return errors.New("device is already attached")
	}
	b.attached[deviceID] = true
	return nil
}

func (b *Backend) Detach(deviceID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.attached[deviceID] {
		return errors.New("device is not attached")
	}
	delete(b.attached, deviceID)
	return nil
}

func (b *Backend) exists(deviceID string) bool {
	for _, item := range b.devices {
		if item.ID == deviceID {
			return true
		}
	}
	return false
}
