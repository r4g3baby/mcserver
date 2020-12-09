package packets

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
)

type (
	PacketStatusOutResponse struct {
		Response Response
	}

	Response struct {
		Version     Version     `json:"version"`
		Players     Players     `json:"players"`
		Description Description `json:"description"`
	}

	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	}

	Players struct {
		Max    int      `json:"max"`
		Online int      `json:"online"`
		Sample []Sample `json:"sample"`
	}

	Sample struct {
		Name string    `json:"name"`
		Id   uuid.UUID `json:"id"`
	}

	Description []chat.Component
)

func (packet *PacketStatusOutResponse) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Status, protocol.ClientBound, packet)
}

func (packet *PacketStatusOutResponse) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	response, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}

	var obj Response
	if err := json.Unmarshal([]byte(response), &obj); err != nil {
		return err
	}
	packet.Response = obj

	return nil
}

func (packet *PacketStatusOutResponse) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	response, err := json.Marshal(packet.Response)
	if err != nil {
		return err
	}

	if err := buffer.WriteUtf(string(response), 32767); err != nil {
		return err
	}

	return nil
}

func (d *Description) UnmarshalJSON(data []byte) error {
	desc, err := chat.FromJSON(data)
	if err != nil {
		return err
	}

	*d = desc
	return nil
}

func (d Description) MarshalJSON() ([]byte, error) {
	return chat.ToJSON(d)
}
