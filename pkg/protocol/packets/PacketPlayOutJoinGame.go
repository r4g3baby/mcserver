package packets

import (
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/nbt"
)

type PacketPlayOutJoinGame struct {
	EntityID         int32
	Hardcore         bool
	Gamemode         uint8
	PreviousGamemode int8
	WorldNames       []string
	DimensionCodec   nbt.Tag
	Dimension        nbt.Tag
	WorldName        string
	HashedSeed       int64
	MaxPlayers       int32
	ViewDistance     int32
	ReducedDebug     bool
	RespawnScreen    bool
	IsDebug          bool
	IsFlat           bool
}

var (
	Overworld = nbt.CompoundTag{
		"name": nbt.StringTag("minecraft:overworld"),
		"id":   nbt.IntTag(0),
		"element": nbt.CompoundTag{
			"bed_works":            nbt.ByteTag(1),
			"has_ceiling":          nbt.ByteTag(0),
			"coordinate_scale":     nbt.DoubleTag(1),
			"piglin_safe":          nbt.ByteTag(0),
			"has_skylight":         nbt.ByteTag(1),
			"ultrawarm":            nbt.ByteTag(0),
			"infiniburn":           nbt.StringTag("minecraft:infiniburn_overworld"),
			"effects":              nbt.StringTag("minecraft:overworld"),
			"has_raids":            nbt.ByteTag(1),
			"ambient_light":        nbt.FloatTag(0),
			"logical_height":       nbt.IntTag(256),
			"natural":              nbt.ByteTag(1),
			"respawn_anchor_works": nbt.ByteTag(0),
		},
	}
	OverworldCaves = nbt.CompoundTag{
		"name": nbt.StringTag("minecraft:overworld_caves"),
		"id":   nbt.IntTag(1),
		"element": nbt.CompoundTag{
			"bed_works":            nbt.ByteTag(1),
			"has_ceiling":          nbt.ByteTag(1),
			"coordinate_scale":     nbt.DoubleTag(1),
			"piglin_safe":          nbt.ByteTag(0),
			"has_skylight":         nbt.ByteTag(1),
			"ultrawarm":            nbt.ByteTag(0),
			"infiniburn":           nbt.StringTag("minecraft:infiniburn_overworld"),
			"effects":              nbt.StringTag("minecraft:overworld"),
			"has_raids":            nbt.ByteTag(1),
			"ambient_light":        nbt.FloatTag(0),
			"logical_height":       nbt.IntTag(256),
			"natural":              nbt.ByteTag(1),
			"respawn_anchor_works": nbt.ByteTag(0),
		},
	}
	TheNether = nbt.CompoundTag{
		"name": nbt.StringTag("minecraft:the_nether"),
		"id":   nbt.IntTag(2),
		"element": nbt.CompoundTag{
			"bed_works":            nbt.ByteTag(0),
			"has_ceiling":          nbt.ByteTag(1),
			"coordinate_scale":     nbt.DoubleTag(8),
			"piglin_safe":          nbt.ByteTag(1),
			"has_skylight":         nbt.ByteTag(0),
			"ultrawarm":            nbt.ByteTag(1),
			"infiniburn":           nbt.StringTag("minecraft:infiniburn_nether"),
			"effects":              nbt.StringTag("minecraft:the_nether"),
			"has_raids":            nbt.ByteTag(0),
			"ambient_light":        nbt.FloatTag(0.1),
			"logical_height":       nbt.IntTag(128),
			"natural":              nbt.ByteTag(0),
			"respawn_anchor_works": nbt.ByteTag(1),
			"fixed_time":           nbt.LongTag(18000),
		},
	}
	TheEnd = nbt.CompoundTag{
		"name": nbt.StringTag("minecraft:the_end"),
		"id":   nbt.IntTag(3),
		"element": nbt.CompoundTag{
			"bed_works":            nbt.ByteTag(0),
			"has_ceiling":          nbt.ByteTag(0),
			"coordinate_scale":     nbt.DoubleTag(1),
			"piglin_safe":          nbt.ByteTag(0),
			"has_skylight":         nbt.ByteTag(0),
			"ultrawarm":            nbt.ByteTag(0),
			"infiniburn":           nbt.StringTag("minecraft:infiniburn_end"),
			"effects":              nbt.StringTag("minecraft:the_end"),
			"has_raids":            nbt.ByteTag(1),
			"ambient_light":        nbt.FloatTag(0),
			"logical_height":       nbt.IntTag(256),
			"natural":              nbt.ByteTag(0),
			"respawn_anchor_works": nbt.ByteTag(0),
			"fixed_time":           nbt.LongTag(6000),
		},
	}
	DimensionCodec = nbt.CompoundTag{
		"minecraft:dimension_type": nbt.CompoundTag{
			"type": nbt.StringTag("minecraft:dimension_type"),
			"value": nbt.ListTag{
				Overworld,
				OverworldCaves,
				TheNether,
				TheEnd,
			},
		},
		"minecraft:worldgen/biome": nbt.CompoundTag{
			"type": nbt.StringTag("minecraft:worldgen/biome"),
			"value": nbt.ListTag{
				nbt.CompoundTag{
					"name": nbt.StringTag("minecraft:plains"),
					"id":   nbt.IntTag(1),
					"element": nbt.CompoundTag{
						"category":    nbt.StringTag("plains"),
						"temperature": nbt.FloatTag(0.800000011920929),
						"downfall":    nbt.FloatTag(0.4000000059604645),
						"depth":       nbt.FloatTag(0.125),
						"effects": nbt.CompoundTag{
							"water_fog_color": nbt.IntTag(329011),
							"water_color":     nbt.IntTag(4159204),
							"fog_color":       nbt.IntTag(12638463),
							"mood_sound": nbt.CompoundTag{
								"offset":              nbt.DoubleTag(2),
								"sound":               nbt.StringTag("minecraft:ambient.cave"),
								"block_search_extent": nbt.IntTag(8),
								"tick_delay":          nbt.IntTag(6000),
							},
							"sky_color": nbt.IntTag(7907327),
						},
						"precipitation": nbt.StringTag("rain"),
						"scale":         nbt.FloatTag(0.05000000074505806),
					},
				},
			},
		},
	}
)

func (packet *PacketPlayOutJoinGame) GetID() int32 {
	return 0x24
}

func (packet *PacketPlayOutJoinGame) Read(buffer *bytes.Buffer) error {
	entityID, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	packet.EntityID = entityID

	hardcore, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.Hardcore = hardcore

	gamemode, err := buffer.ReadUint8()
	if err != nil {
		return err
	}
	packet.Gamemode = gamemode

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

	_, dimensionCodec, err := nbt.Read(buffer)
	if err != nil {
		return err
	}
	packet.DimensionCodec = dimensionCodec

	_, dimension, err := nbt.Read(buffer)
	if err != nil {
		return err
	}
	packet.Dimension = dimension

	worldName, err := buffer.ReadUtf(32767)
	if err != nil {
		return err
	}
	packet.WorldName = worldName

	hashedSeed, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	packet.HashedSeed = hashedSeed

	maxPlayers, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.MaxPlayers = maxPlayers

	viewDistance, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.ViewDistance = viewDistance

	reducedDebug, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.ReducedDebug = reducedDebug

	respawnScreen, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.RespawnScreen = respawnScreen

	debug, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.IsDebug = debug

	flat, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.IsFlat = flat

	return nil
}

func (packet *PacketPlayOutJoinGame) Write(buffer *bytes.Buffer) error {
	if err := buffer.WriteInt32(packet.EntityID); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.Hardcore); err != nil {
		return err
	}

	if err := buffer.WriteUint8(packet.Gamemode); err != nil {
		return err
	}

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

	if err := nbt.Write(buffer, "", packet.DimensionCodec); err != nil {
		return err
	}

	if err := nbt.Write(buffer, "", packet.Dimension); err != nil {
		return err
	}

	if err := buffer.WriteUtf(packet.WorldName, 32767); err != nil {
		return err
	}

	if err := buffer.WriteInt64(packet.HashedSeed); err != nil {
		return err
	}

	if err := buffer.WriteVarInt(packet.MaxPlayers); err != nil {
		return err
	}

	if err := buffer.WriteVarInt(packet.ViewDistance); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.ReducedDebug); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.RespawnScreen); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.IsDebug); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.IsFlat); err != nil {
		return err
	}

	return nil
}
