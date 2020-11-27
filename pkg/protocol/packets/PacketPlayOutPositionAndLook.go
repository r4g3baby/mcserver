package packets

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type PacketPlayOutPositionAndLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      uint8
	TeleportID int32
}

func (packet *PacketPlayOutPositionAndLook) GetID() int32 {
	return 0x34
}

func (packet *PacketPlayOutPositionAndLook) Read(buffer *bytes.Buffer) error {
	x, err := buffer.ReadFloat64()
	if err != nil {
		return err
	}
	packet.X = x

	y, err := buffer.ReadFloat64()
	if err != nil {
		return err
	}
	packet.Y = y

	z, err := buffer.ReadFloat64()
	if err != nil {
		return err
	}
	packet.Z = z

	yaw, err := buffer.ReadFloat32()
	if err != nil {
		return err
	}
	packet.Yaw = yaw

	pitch, err := buffer.ReadFloat32()
	if err != nil {
		return err
	}
	packet.Pitch = pitch

	flags, err := buffer.ReadUint8()
	if err != nil {
		return err
	}
	packet.Flags = flags

	teleportID, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.TeleportID = teleportID

	return nil
}

func (packet *PacketPlayOutPositionAndLook) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteFloat64(packet.X); err != nil {
		return err
	}

	if err := buffer.WriteFloat64(packet.Y); err != nil {
		return err
	}

	if err := buffer.WriteFloat64(packet.Z); err != nil {
		return err
	}

	if err := buffer.WriteFloat32(packet.Yaw); err != nil {
		return err
	}

	if err := buffer.WriteFloat32(packet.Pitch); err != nil {
		return err
	}

	if err := buffer.WriteUint8(packet.Flags); err != nil {
		return err
	}

	if err := buffer.WriteVarInt(packet.TeleportID); err != nil {
		return err
	}

	return nil
}
