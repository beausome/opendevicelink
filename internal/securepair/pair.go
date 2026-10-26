// Package securepair implements local pairing primitives without external dependencies.
package securepair

import (
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

type SessionKey [32]byte

type Peer struct {
	private *ecdh.PrivateKey
}

func NewPeer() (*Peer, error) {
	private, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Peer{private: private}, nil
}

func (p *Peer) PublicKey() []byte {
	return append([]byte(nil), p.private.PublicKey().Bytes()...)
}

func (p *Peer) Derive(remotePublic []byte, context string) (SessionKey, error) {
	if context == "" {
		return SessionKey{}, errors.New("key context is required")
	}
	remote, err := ecdh.X25519().NewPublicKey(remotePublic)
	if err != nil {
		return SessionKey{}, fmt.Errorf("parse remote public key: %w", err)
	}
	shared, err := p.private.ECDH(remote)
	if err != nil {
		return SessionKey{}, fmt.Errorf("derive shared secret: %w", err)
	}
	mac := hmac.New(sha256.New, shared)
	mac.Write([]byte("OpenDeviceLink/session/v1"))
	mac.Write([]byte{0})
	mac.Write([]byte(context))
	var key SessionKey
	copy(key[:], mac.Sum(nil))
	return key, nil
}

func GenerateCode() (string, error) {
	var source [4]byte
	if _, err := rand.Read(source[:]); err != nil {
		return "", err
	}
	value := binary.BigEndian.Uint32(source[:]) % 1_000_000
	return fmt.Sprintf("%06d", value), nil
}

func ValidCode(value string) bool {
	if len(value) != 6 {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func VerificationCode(key SessionKey) string {
	value := binary.BigEndian.Uint32(key[:4]) % 1_000_000
	return fmt.Sprintf("%06d", value)
}
