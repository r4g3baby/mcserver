package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketLoginOutCompression struct {
	Threshold int32
}

func (packet *PacketLoginOutCompression) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Login, protocol.ClientBound, packet)
}

func (packet *PacketLoginOutCompression) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	threshold, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.Threshold = threshold

	return nil
}

func (packet *PacketLoginOutCompression) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteVarInt(packet.Threshold); err != nil {
		return err
	}

	return nil
}
