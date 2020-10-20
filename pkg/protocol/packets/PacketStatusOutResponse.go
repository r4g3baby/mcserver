package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketStatusOutResponse struct {
	Response string
}

func (packet *PacketStatusOutResponse) GetID() int32 {
	return 0x00
}

func (packet *PacketStatusOutResponse) Read(buffer *bytes.Buffer) error {
	response, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}
	packet.Response = response

	return nil
}

func (packet *PacketStatusOutResponse) Write(buffer *bytes.Buffer) error {
	err := buffer.WriteUtf(packet.Response, 32767)
	if err != nil {
		return err
	}

	return nil
}
