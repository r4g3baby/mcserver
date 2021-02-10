package packets

import (
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
)

type PacketPlayOutChatMessage struct {
	Message  []chat.Component
	Position int8
	Sender   uuid.UUID
}

func (packet *PacketPlayOutChatMessage) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Play, protocol.ClientBound, packet)
}

func (packet *PacketPlayOutChatMessage) Read(proto protocol.Protocol, buffer *bytes.Buffer) error {
	messageStr, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}

	message, err := chat.FromJSON([]byte(messageStr))
	if err != nil {
		return err
	}
	packet.Message = message

	position, err := buffer.ReadInt8()
	if err != nil {
		return err
	}
	packet.Position = position

	if proto >= protocol.V1_16 {
		sender, err := buffer.ReadUUID()
		if err != nil {
			return err
		}
		packet.Sender = sender
	}

	return nil
}

func (packet *PacketPlayOutChatMessage) Write(proto protocol.Protocol, buffer *bytes.Buffer) error {
	message, err := chat.ToJSON(packet.Message)
	if err != nil {
		return err
	}

	if err := buffer.WriteUtf(string(message), 32767); err != nil {
		return err
	}

	if err := buffer.WriteInt8(packet.Position); err != nil {
		return err
	}

	if proto >= protocol.V1_16 {
		if err := buffer.WriteUUID(packet.Sender); err != nil {
			return err
		}
	}

	return nil
}
