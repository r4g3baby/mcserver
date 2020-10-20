package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketStatusOutPong struct {
	Payload int64
}

func (packet *PacketStatusOutPong) GetID() int32 {
	return 0x01
}

func (packet *PacketStatusOutPong) Read(buffer *bytes.Buffer) error {
	payload, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.Payload = payload

	return nil
}

func (packet *PacketStatusOutPong) Write(buffer *bytes.Buffer) error {
	err := buffer.WriteInt64(packet.Payload)
	if err != nil {
		return err
	}

	return nil
}
