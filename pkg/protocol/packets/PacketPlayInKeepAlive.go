package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketPlayInKeepAlive struct {
	KeepAliveID int32
}

func (packet *PacketPlayInKeepAlive) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Play, protocol.ServerBound, packet)
}

func (packet *PacketPlayInKeepAlive) Read(proto protocol.Protocol, buffer *bytes.Buffer) error {
	if proto >= protocol.V1_12_2 {
		keepAliveID, err := buffer.ReadInt64()
		if err != nil {
			return err
		}
		packet.KeepAliveID = int32(keepAliveID)
	} else {
		keepAliveID, err := buffer.ReadVarInt()
		if err != nil {
			return err
		}
		packet.KeepAliveID = keepAliveID
	}

	return nil
}

func (packet *PacketPlayInKeepAlive) Write(proto protocol.Protocol, buffer *bytes.Buffer) error {
	if proto >= protocol.V1_12_2 {
		if err := buffer.WriteInt64(int64(packet.KeepAliveID)); err != nil {
			return err
		}
	} else {
		if err := buffer.WriteVarInt(packet.KeepAliveID); err != nil {
			return err
		}
	}

	return nil
}
