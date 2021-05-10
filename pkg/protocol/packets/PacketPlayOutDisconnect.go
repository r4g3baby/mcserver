package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
)

type PacketPlayOutDisconnect struct {
	Reason []chat.Component
}

func (packet *PacketPlayOutDisconnect) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Play, protocol.ClientBound, packet)
}

func (packet *PacketPlayOutDisconnect) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	reasonStr, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}

	reason, err := chat.FromJSON([]byte(reasonStr))
	if err != nil {
		return err
	}
	packet.Reason = reason

	return nil
}

func (packet *PacketPlayOutDisconnect) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	reason, err := chat.ToJSON(packet.Reason)
	if err != nil {
		return err
	}

	if err := buffer.WriteUtf(string(reason), 32767); err != nil {
		return err
	}

	return nil
}
