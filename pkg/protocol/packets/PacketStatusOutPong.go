package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketStatusOutPong struct {
	Payload int64
}

func (packet *PacketStatusOutPong) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Status, protocol.ClientBound, packet)
}

func (packet *PacketStatusOutPong) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	payload, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.Payload = payload

	return nil
}

func (packet *PacketStatusOutPong) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteInt64(packet.Payload); err != nil {
		return err
	}

	return nil
}
