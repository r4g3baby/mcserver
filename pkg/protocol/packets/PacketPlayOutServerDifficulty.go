package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketPlayOutServerDifficulty struct {
	Difficulty uint8
	Locked     bool
}

func (packet *PacketPlayOutServerDifficulty) GetID() int32 {
	return 0x0D
}

func (packet *PacketPlayOutServerDifficulty) Read(buffer *bytes.Buffer) error {
	difficulty, err := buffer.ReadUint8()
	if err != nil {
		return err
	}
	packet.Difficulty = difficulty

	locked, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.Locked = locked

	return nil
}

func (packet *PacketPlayOutServerDifficulty) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteUint8(packet.Difficulty); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.Locked); err != nil {
		return err
	}

	return nil
}
