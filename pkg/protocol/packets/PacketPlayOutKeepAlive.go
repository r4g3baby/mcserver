package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketPlayOutKeepAlive struct {
	KeepAliveID int64
}

func (packet *PacketPlayOutKeepAlive) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Play, protocol.ClientBound, packet)
}

func (packet *PacketPlayOutKeepAlive) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	keepAliveID, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.KeepAliveID = keepAliveID

	return nil
}

func (packet *PacketPlayOutKeepAlive) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	err := buffer.WriteInt64(packet.KeepAliveID)
	if err != nil {
		return err
	}

	return nil
}
