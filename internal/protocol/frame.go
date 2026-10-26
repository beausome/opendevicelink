// Package protocol defines bounded binary frames for future device transport.
package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const MaxPayloadSize = 1 << 20

type Type byte

const (
	TypeControl Type = 1
	TypeData    Type = 2
	TypeAck     Type = 3
)

type Frame struct {
	Type    Type
	Payload []byte
}

func Encode(writer io.Writer, frame Frame) error {
	if !validType(frame.Type) {
		return fmt.Errorf("invalid frame type %d", frame.Type)
	}
	if len(frame.Payload) > MaxPayloadSize {
		return fmt.Errorf("payload is %d bytes; maximum is %d", len(frame.Payload), MaxPayloadSize)
	}
	length := uint32(len(frame.Payload) + 1)
	if err := binary.Write(writer, binary.BigEndian, length); err != nil {
		return fmt.Errorf("write frame length: %w", err)
	}
	if _, err := writer.Write([]byte{byte(frame.Type)}); err != nil {
		return fmt.Errorf("write frame type: %w", err)
	}
	if _, err := writer.Write(frame.Payload); err != nil {
		return fmt.Errorf("write frame payload: %w", err)
	}
	return nil
}

func Decode(reader io.Reader) (Frame, error) {
	var length uint32
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return Frame{}, fmt.Errorf("read frame length: %w", err)
	}
	if length == 0 {
		return Frame{}, errors.New("frame length must include a type byte")
	}
	if length > MaxPayloadSize+1 {
		return Frame{}, fmt.Errorf("frame length %d exceeds maximum %d", length, MaxPayloadSize+1)
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return Frame{}, fmt.Errorf("read frame data: %w", err)
	}
	typeValue := Type(data[0])
	if !validType(typeValue) {
		return Frame{}, fmt.Errorf("invalid frame type %d", typeValue)
	}
	return Frame{Type: typeValue, Payload: data[1:]}, nil
}

func validType(value Type) bool {
	return value == TypeControl || value == TypeData || value == TypeAck
}
