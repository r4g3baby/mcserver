package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/nbt"
)

type (
	PacketPlayOutJoinGame struct {
		EntityID         int32
		Hardcore         bool
		Gamemode         uint8
		PreviousGamemode int8
		WorldNames       []string
		DimensionCodec   protocol.DimensionCodec
		Dimension        protocol.Dimension
		WorldName        string
		DimensionID      int8
		Difficulty       uint8
		HashedSeed       int64
		MaxPlayers       int32
		LevelType        string
		ViewDistance     int32
		ReducedDebug     bool
		RespawnScreen    bool
		IsDebug          bool
		IsFlat           bool
	}
)

func (packet *PacketPlayOutJoinGame) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Play, protocol.ClientBound, packet)
}

func (packet *PacketPlayOutJoinGame) Read(proto protocol.Protocol, buffer *bytes.Buffer) error {
	entityID, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	packet.EntityID = entityID

	if proto >= protocol.V1_16_2 {
		hardcore, err := buffer.ReadBool()
		if err != nil {
			return err
		}
		packet.Hardcore = hardcore
	}

	gamemode, err := buffer.ReadUint8()
	if err != nil {
		return err
	}
	packet.Gamemode = gamemode

	if proto >= protocol.V1_16 {
		previousGamemode, err := buffer.ReadInt8()
		if err != nil {
			return err
		}
		packet.PreviousGamemode = previousGamemode

		worldCount, err := buffer.ReadVarInt()
		if err != nil {
			return err
		}

		var worldNames []string
		for i := worldCount; i > 0; i-- {
			worldName, err := buffer.ReadUtf(32767)
			if err != nil {
				return err
			}
			worldNames = append(worldNames, worldName)
		}
		packet.WorldNames = worldNames

		_, dimensionCodecTag, err := nbt.Read(buffer)
		if err != nil {
			return err
		}
		dimensionCodec, err := protocol.DimensionCodecFromTag(dimensionCodecTag, proto)
		if err != nil {
			return err
		}
		packet.DimensionCodec = dimensionCodec

		if proto >= protocol.V1_16_2 {
			_, dimensionTag, err := nbt.Read(buffer)
			if err != nil {
				return err
			}

			dimension, err := protocol.DimensionFromTag(dimensionTag)
			if err != nil {
				return err
			}

			packet.Dimension = dimension
			for _, dim := range packet.DimensionCodec.Dimensions {
				if dim.ID == packet.Dimension.ID {
					packet.Dimension.Name = dim.Name
					break
				}
			}
		} else {
			dimensionName, err := buffer.ReadUtf(32767)
			if err != nil {
				return err
			}

			packet.Dimension = protocol.Dimension{Name: dimensionName}
			for _, dim := range packet.DimensionCodec.Dimensions {
				if dim.Name == dimensionName {
					packet.Dimension = dim
					break
				}
			}
		}

		worldName, err := buffer.ReadUtf(32767)
		if err != nil {
			return err
		}
		packet.WorldName = worldName
	} else {
		if proto >= protocol.V1_9_1 {
			dimensionID, err := buffer.ReadInt32()
			if err != nil {
				return err
			}
			packet.DimensionID = int8(dimensionID)
		} else {
			dimensionID, err := buffer.ReadInt8()
			if err != nil {
				return err
			}
			packet.DimensionID = dimensionID
		}
	}

	if proto < protocol.V1_14 {
		difficulty, err := buffer.ReadUint8()
		if err != nil {
			return err
		}
		packet.Difficulty = difficulty
	}

	if proto >= protocol.V1_15 {
		hashedSeed, err := buffer.ReadInt64()
		if err != nil {
			return err
		}
		packet.HashedSeed = hashedSeed
	}

	if proto >= protocol.V1_16 {
		maxPlayers, err := buffer.ReadVarInt()
		if err != nil {
			return err
		}
		packet.MaxPlayers = maxPlayers
	} else {
		maxPlayers, err := buffer.ReadUint8()
		if err != nil {
			return err
		}
		packet.MaxPlayers = int32(maxPlayers)

		levelType, err := buffer.ReadUtf(16)
		if err != nil {
			return err
		}
		packet.LevelType = levelType
	}

	if proto >= protocol.V1_14 {
		viewDistance, err := buffer.ReadVarInt()
		if err != nil {
			return err
		}
		packet.ViewDistance = viewDistance
	}

	reducedDebug, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.ReducedDebug = reducedDebug

	if proto >= protocol.V1_15 {
		respawnScreen, err := buffer.ReadBool()
		if err != nil {
			return err
		}
		packet.RespawnScreen = respawnScreen
	}

	if proto >= protocol.V1_16 {
		isDebug, err := buffer.ReadBool()
		if err != nil {
			return err
		}
		packet.IsDebug = isDebug

		isFlat, err := buffer.ReadBool()
		if err != nil {
			return err
		}
		packet.IsFlat = isFlat
	}

	return nil
}

func (packet *PacketPlayOutJoinGame) Write(proto protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteInt32(packet.EntityID); err != nil {
		return err
	}

	if proto >= protocol.V1_16_2 {
		if err := buffer.WriteBool(packet.Hardcore); err != nil {
			return err
		}
	}

	if err := buffer.WriteUint8(packet.Gamemode); err != nil {
		return err
	}

	if proto >= protocol.V1_16 {
		if err := buffer.WriteInt8(packet.PreviousGamemode); err != nil {
			return err
		}

		if err := buffer.WriteVarInt(int32(len(packet.WorldNames))); err != nil {
			return err
		}

		for _, worldName := range packet.WorldNames {
			if err := buffer.WriteUtf(worldName, 32767); err != nil {
				return err
			}
		}

		if err := nbt.Write(buffer, "", packet.DimensionCodec.ToCompound(proto)); err != nil {
			return err
		}

		if proto >= protocol.V1_16_2 {
			if err := nbt.Write(buffer, "", packet.Dimension.ToCompound(proto)); err != nil {
				return err
			}
		} else {
			if err := buffer.WriteUtf(packet.Dimension.Name, 32767); err != nil {
				return err
			}
		}

		if err := buffer.WriteUtf(packet.WorldName, 32767); err != nil {
			return err
		}
	} else {
		if proto >= protocol.V1_9_1 {
			if err := buffer.WriteInt32(int32(packet.DimensionID)); err != nil {
				return err
			}
		} else {
			if err := buffer.WriteInt8(packet.DimensionID); err != nil {
				return err
			}
		}
	}

	if proto < protocol.V1_14 {
		if err := buffer.WriteUint8(packet.Difficulty); err != nil {
			return err
		}
	}

	if proto >= protocol.V1_15 {
		if err := buffer.WriteInt64(packet.HashedSeed); err != nil {
			return err
		}
	}

	if proto >= protocol.V1_16 {
		if err := buffer.WriteVarInt(packet.MaxPlayers); err != nil {
			return err
		}
	} else {
		if err := buffer.WriteUint8(uint8(packet.MaxPlayers)); err != nil {
			return err
		}

		if err := buffer.WriteUtf(packet.LevelType, 16); err != nil {
			return err
		}
	}

	if proto >= protocol.V1_14 {
		if err := buffer.WriteVarInt(packet.ViewDistance); err != nil {
			return err
		}
	}

	if err := buffer.WriteBool(packet.ReducedDebug); err != nil {
		return err
	}

	if proto >= protocol.V1_15 {
		if err := buffer.WriteBool(packet.RespawnScreen); err != nil {
			return err
		}
	}

	if proto >= protocol.V1_16 {
		if err := buffer.WriteBool(packet.IsDebug); err != nil {
			return err
		}

		if err := buffer.WriteBool(packet.IsFlat); err != nil {
			return err
		}
	}

	return nil
}
