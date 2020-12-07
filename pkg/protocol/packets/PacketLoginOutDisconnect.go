package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketLoginOutDisconnect struct {
	Reason string
}

func (packet *PacketLoginOutDisconnect) GetID() int32 {
	return 0x00
}

func (packet *PacketLoginOutDisconnect) Read(buffer *bytes.Buffer) error {
	reason, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}
	packet.Reason = reason

	return nil
}

func (packet *PacketLoginOutDisconnect) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteUtf(packet.Reason, 32767); err != nil {
		return err
	}

	return nil
}
