package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
)

type PacketPlayOutPositionAndLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      uint8
	TeleportID int32
}

func (packet *PacketPlayOutPositionAndLook) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Play, protocol.ClientBound, packet)
}

func (packet *PacketPlayOutPositionAndLook) Read(proto protocol.Protocol, buffer *bytes.Buffer) error {
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

	if proto >= protocol.V1_9 {
		teleportID, err := buffer.ReadVarInt()
		if err != nil {
			return err
		}
		packet.TeleportID = teleportID
	}

	return nil
}

func (packet *PacketPlayOutPositionAndLook) Write(proto protocol.Protocol, buffer *bytes.Buffer) error {
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

	if proto >= protocol.V1_9 {
		if err := buffer.WriteVarInt(packet.TeleportID); err != nil {
			return err
		}
	}

	return nil
}
