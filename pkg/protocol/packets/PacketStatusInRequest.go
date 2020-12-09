package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketStatusInRequest struct{}

func (packet *PacketStatusInRequest) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Status, protocol.ServerBound, packet)
}

func (packet *PacketStatusInRequest) Read(_ protocol.Protocol, _ *bytes.Buffer) error {
	return nil
}

func (packet *PacketStatusInRequest) Write(_ protocol.Protocol, _ *bytes.Buffer) error {
	return nil
}
