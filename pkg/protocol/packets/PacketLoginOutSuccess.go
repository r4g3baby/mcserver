package packets

import (
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketLoginOutSuccess struct {
	UniqueID uuid.UUID
	Username string
}

func (packet *PacketLoginOutSuccess) GetID() int32 {
	return 0x02
}

func (packet *PacketLoginOutSuccess) Read(buffer *bytes.Buffer) error {
	uniqueID, err := buffer.ReadUUID()
	if err != nil {
		return err
	}
	packet.UniqueID = uniqueID

	username, err := buffer.ReadUtf(16)
	if err != nil {
		return err
	}
	packet.Username = username

	return nil
}

func (packet *PacketLoginOutSuccess) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteUUID(packet.UniqueID); err != nil {
		return err
	}

	if err := buffer.WriteUtf(packet.Username, 16); err != nil {
		return err
	}

	return nil
}
