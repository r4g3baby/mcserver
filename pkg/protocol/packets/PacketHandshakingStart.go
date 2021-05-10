package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketHandshakingStart struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (packet *PacketHandshakingStart) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Handshaking, protocol.ServerBound, packet)
}

func (packet *PacketHandshakingStart) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	protocolVersion, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.ProtocolVersion = protocolVersion

	serverAddress, err := buffer.ReadUtf(255)
	if err != nil {
		return err
	}
	packet.ServerAddress = serverAddress

	serverPort, err := buffer.ReadUint16()
	if err != nil {
		return err
	}
	packet.ServerPort = serverPort

	nextState, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.NextState = nextState

	return nil
}

func (packet *PacketHandshakingStart) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteVarInt(packet.ProtocolVersion); err != nil {
		return err
	}

	if err := buffer.WriteUtf(packet.ServerAddress, 255); err != nil {
		return err
	}

	if err := buffer.WriteUint16(packet.ServerPort); err != nil {
		return err
	}

	if err := buffer.WriteVarInt(packet.NextState); err != nil {
		return err
	}

	return nil
}
