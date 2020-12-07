package packets

import (
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
)

type PacketPlayOutDisconnect struct {
	Reason []chat.Component
}

func (packet *PacketPlayOutDisconnect) GetID() int32 {
	return 0x19
}

func (packet *PacketPlayOutDisconnect) Read(buffer *bytes.Buffer) error {
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

func (packet *PacketPlayOutDisconnect) Write(buffer *bytes.Buffer) error {
	reason, err := chat.ToJSON(packet.Reason)
	if err != nil {
		return err
	}

	if err := buffer.WriteUtf(string(reason), 32767); err != nil {
		return err
	}

	return nil
}
