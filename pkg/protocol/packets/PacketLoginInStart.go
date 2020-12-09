package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketLoginInStart struct {
	Username string
}

func (packet *PacketLoginInStart) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Login, protocol.ServerBound, packet)
}

func (packet *PacketLoginInStart) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	username, err := buffer.ReadUtf(16)
	if err != nil {
		return err
	}
	packet.Username = username

	return nil
}

func (packet *PacketLoginInStart) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteUtf(packet.Username, 16); err != nil {
		return err
	}

	return nil
}
