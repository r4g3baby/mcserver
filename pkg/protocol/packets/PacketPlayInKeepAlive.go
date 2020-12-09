package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketPlayInKeepAlive struct {
	KeepAliveID int64
}

func (packet *PacketPlayInKeepAlive) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Play, protocol.ServerBound, packet)
}

func (packet *PacketPlayInKeepAlive) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	keepAliveID, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.KeepAliveID = keepAliveID

	return nil
}

func (packet *PacketPlayInKeepAlive) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	err := buffer.WriteInt64(packet.KeepAliveID)
	if err != nil {
		return err
	}

	return nil
}
