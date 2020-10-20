package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketStatusInRequest struct{}

func (packet *PacketStatusInRequest) GetID() int32 {
	return 0x00
}

func (packet *PacketStatusInRequest) Read(_ *bytes.Buffer) error {
	return nil
}

func (packet *PacketStatusInRequest) Write(_ *bytes.Buffer) error {
	return nil
}
