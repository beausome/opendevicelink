// Package lease provides exclusive, expiring device ownership.
package lease

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrBusy     = errors.New("device already has an active lease")
	ErrNotFound = errors.New("lease not found")
)

type Lease struct {
	DeviceID string
	Holder   string
	Token    string
	Expires  time.Time
}

type Manager struct {
	mu       sync.Mutex
	now      func() time.Time
	byDevice map[string]Lease
	byToken  map[string]string
}

func NewManager(now func() time.Time) *Manager {
	return &Manager{
		now:      now,
		byDevice: make(map[string]Lease),
		byToken:  make(map[string]string),
	}
}

func (m *Manager) Acquire(deviceID, holder string, ttl time.Duration) (Lease, error) {
	if deviceID == "" || holder == "" {
		return Lease{}, errors.New("device ID and holder are required")
	}
	if ttl <= 0 {
		return Lease{}, errors.New("lease duration must be positive")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()
	if current, exists := m.byDevice[deviceID]; exists {
		if now.Before(current.Expires) {
			return Lease{}, fmt.Errorf("%w: held by %s until %s", ErrBusy, current.Holder, current.Expires.UTC().Format(time.RFC3339))
		}
		delete(m.byToken, current.Token)
	}

	token, err := randomToken()
	if err != nil {
		return Lease{}, fmt.Errorf("create lease token: %w", err)
	}
	created := Lease{DeviceID: deviceID, Holder: holder, Token: token, Expires: now.Add(ttl)}
	m.byDevice[deviceID] = created
	m.byToken[token] = deviceID
	return created, nil
}

func (m *Manager) Release(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	deviceID, exists := m.byToken[token]
	if !exists {
		return ErrNotFound
	}
	delete(m.byToken, token)
	delete(m.byDevice, deviceID)
	return nil
}

func (m *Manager) Active(deviceID string) (Lease, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, exists := m.byDevice[deviceID]
	if !exists {
		return Lease{}, false
	}
	if !m.now().Before(current.Expires) {
		delete(m.byDevice, deviceID)
		delete(m.byToken, current.Token)
		return Lease{}, false
	}
	return current, true
}

func randomToken() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
