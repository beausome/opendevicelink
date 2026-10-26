package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestFrameRoundTrip(t *testing.T) {
	original := Frame{Type: TypeData, Payload: []byte("usb-packet")}
	var encoded bytes.Buffer
	if err := Encode(&encoded, original); err != nil {
		t.Fatalf("encode: %v", err)
	}
	decoded, err := Decode(&encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.Type != original.Type || !bytes.Equal(decoded.Payload, original.Payload) {
		t.Fatalf("decoded frame = %#v, want %#v", decoded, original)
	}
}

func TestDecodeRejectsOversizedFrameBeforeAllocation(t *testing.T) {
	var input bytes.Buffer
	if err := binary.Write(&input, binary.BigEndian, uint32(MaxPayloadSize+2)); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	if _, err := Decode(&input); err == nil {
		t.Fatal("oversized frame accepted")
	}
}

func TestDecodeRejectsUnknownType(t *testing.T) {
	input := bytes.NewBuffer([]byte{0, 0, 0, 1, 99})
	if _, err := Decode(input); err == nil {
		t.Fatal("unknown frame type accepted")
	}
}
