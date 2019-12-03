package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PacketType is the value of the packet type field in an RCON packet.
type PacketType int32

const (
	// Response is a response to an RCON command.
	Response PacketType = 0

	// Command is an incoming RCON command.
	Command PacketType = 2

	// AuthResponse is a response to an AuthRequest. If the request is rejected,
	// the packet must have a request ID of -1.
	AuthResponse PacketType = 2

	// AuthRequest is an incoming RCON authentication request.
	AuthRequest PacketType = 3
)

// String maps PacketType values to their string names.
func (pt PacketType) String() string {
	switch pt {
	case Response:
		return "Response"
	case Command:
		return "Command"
	case AuthRequest:
		return "AuthRequest"
	default:
		return ""
	}
}

// Packet is an RCON message packet.
type Packet struct {
	// Length is the length in bytes of the RCON packet, minus the length field
	// itself. Its minimum value is 10:
	//
	//	   4 bytes	(RequestID)
	//	+  4 bytes	(Type)
	//	+  0 bytes	(empty Payload)
	//	+  2 bytes	(2 null padding bytes)
	//	= 10 bytes
	Length int32

	// RequestID is the client-provided ID of the packet. Responses must have
	// the same value, except for an AuthResponse to a failed AuthRequest, in
	// which case it must be -1.
	RequestID int32

	Type    PacketType
	Payload []byte
}

// NewPacket builds a packet with calculated length for the given values.
func NewPacket(id int32, pt PacketType, payload []byte) *Packet {
	return &Packet{
		Length:    int32(4 + 4 + len(payload) + 2),
		RequestID: id,
		Type:      pt,
		Payload:   payload,
	}
}

// Bytes is the byte representation of the packet for writing to the network
// connection.
func (packet *Packet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, v := range []interface{}{
		packet.Length,
		packet.RequestID,
		packet.Type,
		packet.Payload,
		[]byte{0, 0},
	} {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// ValidateType ensures that the packet is of the expected type.
func (packet *Packet) ValidateType(pt PacketType) error {
	if packet.Type != pt {
		return fmt.Errorf("Invalid packet type %s, expected %s", packet.Type, pt)
	}

	return nil
}
