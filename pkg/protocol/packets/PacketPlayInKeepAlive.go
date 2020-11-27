package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketPlayInKeepAlive struct {
	KeepAliveID int64
}

func (packet *PacketPlayInKeepAlive) GetID() int32 {
	return 0x10
}

func (packet *PacketPlayInKeepAlive) Read(buffer *bytes.Buffer) error {
	keepAliveID, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.KeepAliveID = keepAliveID

	return nil
}

func (packet *PacketPlayInKeepAlive) Write(buffer *bytes.Buffer) error {
	err := buffer.WriteInt64(packet.KeepAliveID)
	if err != nil {
		return err
	}

	return nil
}
