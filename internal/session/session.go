// Package session enforces the handoff state machine.
package session

import (
	"crypto/subtle"
	"errors"
	"fmt"

	"github.com/beausome/opendevicelink/internal/securepair"
)

type State string

const (
	Discovered State = "discovered"
	Pairing    State = "pairing"
	Authorized State = "authorized"
	Attached   State = "attached"
	Detached   State = "detached"
)

type Session struct {
	deviceID string
	state    State
	code     string
}

func New(deviceID string) *Session {
	return &Session{deviceID: deviceID, state: Discovered}
}

func (s *Session) State() State {
	return s.state
}

func (s *Session) BeginPairing(code string) error {
	if s.deviceID == "" {
		return errors.New("device ID is required")
	}
	if s.state != Discovered {
		return s.transitionError(Pairing)
	}
	if !securepair.ValidCode(code) {
		return errors.New("pairing code must contain exactly six digits")
	}
	s.code = code
	s.state = Pairing
	return nil
}

func (s *Session) Authorize(code string) error {
	if s.state != Pairing {
		return s.transitionError(Authorized)
	}
	if subtle.ConstantTimeCompare([]byte(s.code), []byte(code)) != 1 {
		return errors.New("pairing code does not match")
	}
	s.code = ""
	s.state = Authorized
	return nil
}

func (s *Session) Attach() error {
	if s.state != Authorized {
		return s.transitionError(Attached)
	}
	s.state = Attached
	return nil
}

func (s *Session) Detach() error {
	if s.state != Attached {
		return s.transitionError(Detached)
	}
	s.state = Detached
	return nil
}

func (s *Session) transitionError(target State) error {
	return fmt.Errorf("cannot transition from %s to %s", s.state, target)
}
