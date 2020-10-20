package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketHandshakingStart struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (packet *PacketHandshakingStart) GetID() int32 {
	return 0x00
}

func (packet *PacketHandshakingStart) Read(buffer *bytes.Buffer) error {
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

func (packet *PacketHandshakingStart) Write(buffer *bytes.Buffer) error {
	err := buffer.WriteVarInt(packet.ProtocolVersion)
	if err != nil {
		return err
	}

	err = buffer.WriteUtf(packet.ServerAddress, 255)
	if err != nil {
		return err
	}

	err = buffer.WriteUint16(packet.ServerPort)
	if err != nil {
		return err
	}

	err = buffer.WriteVarInt(packet.NextState)
	if err != nil {
		return err
	}

	return nil
}
