package packets

import (
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
)

type PacketLoginOutDisconnect struct {
	Reason []chat.Component
}

func (packet *PacketLoginOutDisconnect) GetID() int32 {
	return 0x00
}

func (packet *PacketLoginOutDisconnect) Read(buffer *bytes.Buffer) error {
	reasonStr, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}

	reason, err := chat.FromJSON([]byte(reasonStr))
	if err != nil {
		return err
	}
	packet.Reason = reason

	return nil
}

func (packet *PacketLoginOutDisconnect) Write(buffer *bytes.Buffer) error {
	reason, err := chat.ToJSON(packet.Reason)
	if err != nil {
		return err
	}

	if err := buffer.WriteUtf(string(reason), 32767); err != nil {
		return err
	}

	return nil
}
