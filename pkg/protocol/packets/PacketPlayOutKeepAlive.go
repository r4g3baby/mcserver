package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketPlayOutKeepAlive struct {
	KeepAliveID int64
}

func (packet *PacketPlayOutKeepAlive) GetID() int32 {
	return 0x1F
}

func (packet *PacketPlayOutKeepAlive) Read(buffer *bytes.Buffer) error {
	keepAliveID, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.KeepAliveID = keepAliveID

	return nil
}

func (packet *PacketPlayOutKeepAlive) Write(buffer *bytes.Buffer) error {
	err := buffer.WriteInt64(packet.KeepAliveID)
	if err != nil {
		return err
	}

	return nil
}
