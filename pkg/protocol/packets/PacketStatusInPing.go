package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketStatusInPing struct {
	Payload int64
}

func (packet *PacketStatusInPing) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Status, protocol.ServerBound, packet)
}

func (packet *PacketStatusInPing) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	payload, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.Payload = payload

	return nil
}

func (packet *PacketStatusInPing) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteInt64(packet.Payload); err != nil {
		return err
	}

	return nil
}
