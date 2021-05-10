package packets

import (
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketLoginOutSuccess struct {
	UniqueID uuid.UUID
	Username string
}

func (packet *PacketLoginOutSuccess) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Login, protocol.ClientBound, packet)
}

func (packet *PacketLoginOutSuccess) Read(proto protocol.Protocol, buffer *bytes.Buffer) error {
	if proto >= protocol.V1_16 {
		uniqueID, err := buffer.ReadUUID()
		if err != nil {
			return err
		}
		packet.UniqueID = uniqueID
	} else {
		uniqueID, err := buffer.ReadUtf(36)
		if err != nil {
			return err
		}
		parsedUUID, err := uuid.Parse(uniqueID)
		if err != nil {
			return err
		}
		packet.UniqueID = parsedUUID
	}

	username, err := buffer.ReadUtf(16)
	if err != nil {
		return err
	}
	packet.Username = username

	return nil
}

func (packet *PacketLoginOutSuccess) Write(proto protocol.Protocol, buffer *bytes.Buffer) error {
	if proto >= protocol.V1_16 {
		if err := buffer.WriteUUID(packet.UniqueID); err != nil {
			return err
		}
	} else {
		if err := buffer.WriteUtf(packet.UniqueID.String(), 36); err != nil {
			return err
		}
	}

	if err := buffer.WriteUtf(packet.Username, 16); err != nil {
		return err
	}

	return nil
}
