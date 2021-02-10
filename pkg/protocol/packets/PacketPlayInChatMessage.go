package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketPlayInChatMessage struct {
	Message string
}

func (packet *PacketPlayInChatMessage) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Play, protocol.ServerBound, packet)
}

func (packet *PacketPlayInChatMessage) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	message, err := buffer.ReadUtf(256)
	if err != nil {
		return err
	}
	packet.Message = message

	return nil
}

func (packet *PacketPlayInChatMessage) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteUtf(packet.Message, 256); err != nil {
		return err
	}

	return nil
}
