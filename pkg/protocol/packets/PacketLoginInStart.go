package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketLoginInStart struct {
	Username string
}

func (packet *PacketLoginInStart) GetID() int32 {
	return 0x00
}

func (packet *PacketLoginInStart) Read(buffer *bytes.Buffer) error {
	username, err := buffer.ReadUtf(16)
	if err != nil {
		return err
	}
	packet.Username = username

	return nil
}

func (packet *PacketLoginInStart) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteUtf(packet.Username, 16); err != nil {
		return err
	}

	return nil
}
