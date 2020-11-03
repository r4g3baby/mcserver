package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketStatusInPing struct {
	Payload int64
}

func (packet *PacketStatusInPing) GetID() int32 {
	return 0x01
}

func (packet *PacketStatusInPing) Read(buffer *bytes.Buffer) error {
	payload, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.Payload = payload

	return nil
}

func (packet *PacketStatusInPing) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteInt64(packet.Payload); err != nil {
		return err
	}

	return nil
}
