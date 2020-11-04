package packets

import (
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketPlayOutDisconnect struct {
	Reason string
}

func (packet *PacketPlayOutDisconnect) GetID() int32 {
	return 0x19
}

func (packet *PacketPlayOutDisconnect) Read(buffer *bytes.Buffer) error {
	reason, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}
	packet.Reason = reason

	return nil
}

func (packet *PacketPlayOutDisconnect) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteUtf(packet.Reason, 32767); err != nil {
		return err
	}

	return nil
}
