package packets

import (
	"errors"
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
		DimensionCodec   DimensionCodec
		Dimension        Dimension
		WorldName        string
		DimensionID      int32
		HashedSeed       int64
		MaxPlayers       int32
		LevelType        string
		ViewDistance     int32
		ReducedDebug     bool
		RespawnScreen    bool
		IsDebug          bool
		IsFlat           bool
	}

	DimensionCodec struct {
		Dimensions []Dimension
		Biomes     []Biome
	}

	Dimension struct {
		ID                 int32
		Name               string
		BedWorks           bool
		HasCeiling         bool
		CoordinateScale    float64
		PiglinSafe         bool
		HasSkylight        bool
		Ultrawarm          bool
		Infiniburn         string
		Effects            string
		HasRaids           bool
		AmbientLight       float32
		LogicalHeight      int32
		Natural            bool
		RespawnAnchorWorks bool
		FixedTime          *int64
		Shrunk             bool
	}

	Biome struct {
		ID            int32
		Name          string
		Category      string
		Temperature   float32
		Downfall      float32
		Depth         float32
		Precipitation string
		Scale         float32
		Effects       Effects
	}

	Effects struct {
		WaterFogColor int32
		WaterColor    int32
		FogColor      int32
		SkyColor      int32
		MoodSound     MoodSound
	}

	MoodSound struct {
		Offset            float64
		Sound             string
		BlockSearchExtent int32
		TickDelay         int32
	}
)

var (
	Overworld = Dimension{
		ID:                 0,
		Name:               "minecraft:overworld",
		BedWorks:           true,
		HasCeiling:         false,
		CoordinateScale:    1,
		PiglinSafe:         false,
		HasSkylight:        true,
		Ultrawarm:          false,
		Infiniburn:         "minecraft:infiniburn_overworld",
		Effects:            "minecraft:overworld",
		HasRaids:           true,
		AmbientLight:       0,
		LogicalHeight:      256,
		Natural:            true,
		RespawnAnchorWorks: false,
		Shrunk:             true,
	}

	OverworldCaves = Dimension{
		ID:                 1,
		Name:               "minecraft:overworld_caves",
		BedWorks:           true,
		HasCeiling:         true,
		CoordinateScale:    1,
		PiglinSafe:         false,
		HasSkylight:        true,
		Ultrawarm:          false,
		Infiniburn:         "minecraft:infiniburn_overworld",
		Effects:            "minecraft:overworld",
		HasRaids:           true,
		AmbientLight:       0,
		LogicalHeight:      256,
		Natural:            true,
		RespawnAnchorWorks: false,
		Shrunk:             false,
	}

	TheNether = Dimension{
		ID:                 2,
		Name:               "minecraft:the_nether",
		BedWorks:           false,
		HasCeiling:         true,
		CoordinateScale:    8,
		PiglinSafe:         true,
		HasSkylight:        false,
		Ultrawarm:          true,
		Infiniburn:         "minecraft:infiniburn_nether",
		Effects:            "minecraft:the_nether",
		HasRaids:           false,
		AmbientLight:       0.1,
		LogicalHeight:      128,
		Natural:            false,
		RespawnAnchorWorks: true,
		FixedTime:          func(i int64) *int64 { return &i }(18000),
		Shrunk:             true,
	}

	TheEnd = Dimension{
		ID:                 3,
		Name:               "minecraft:the_end",
		BedWorks:           false,
		HasCeiling:         false,
		CoordinateScale:    1,
		PiglinSafe:         false,
		HasSkylight:        false,
		Ultrawarm:          false,
		Infiniburn:         "minecraft:infiniburn_end",
		Effects:            "minecraft:the_end",
		HasRaids:           true,
		AmbientLight:       0,
		LogicalHeight:      256,
		Natural:            false,
		RespawnAnchorWorks: false,
		FixedTime:          func(i int64) *int64 { return &i }(6000),
		Shrunk:             false,
	}

	DefaultDimensionCodec = DimensionCodec{
		Dimensions: []Dimension{Overworld, OverworldCaves, TheNether, TheEnd},
		Biomes: []Biome{{
			ID:            0,
			Name:          "minecraft:ocean",
			Category:      "ocean",
			Temperature:   0.5,
			Downfall:      0.5,
			Depth:         -1,
			Precipitation: "rain",
			Scale:         0.10000000149011612,
			Effects: Effects{
				WaterFogColor: 329011,
				WaterColor:    4159204,
				FogColor:      12638463,
				SkyColor:      8103167,
				MoodSound: MoodSound{
					Offset:            2,
					Sound:             "minecraft:ambient.cave",
					BlockSearchExtent: 8,
					TickDelay:         6000,
				},
			},
		}, {
			ID:            1,
			Name:          "minecraft:plains",
			Category:      "plains",
			Temperature:   0.800000011920929,
			Downfall:      0.4000000059604645,
			Depth:         0.125,
			Precipitation: "rain",
			Scale:         0.05000000074505806,
			Effects: Effects{
				WaterFogColor: 329011,
				WaterColor:    4159204,
				FogColor:      12638463,
				SkyColor:      7907327,
				MoodSound: MoodSound{
					Offset:            2,
					Sound:             "minecraft:ambient.cave",
					BlockSearchExtent: 8,
					TickDelay:         6000,
				},
			},
		}},
	}
)

func (packet *PacketPlayOutJoinGame) GetID(proto protocol.Protocol) (int32, error) {
	return GetPacketID(proto, protocol.Play, protocol.ClientBound, packet)
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
		dimensionCodec, err := DimensionCodecFromTag(dimensionCodecTag, proto)
		if err != nil {
			return err
		}
		packet.DimensionCodec = dimensionCodec

		if proto >= protocol.V1_16_2 {
			_, dimensionTag, err := nbt.Read(buffer)
			if err != nil {
				return err
			}

			dimension, err := DimensionFromTag(dimensionTag)
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

			packet.Dimension = Dimension{Name: dimensionName}
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
		dimensionID, err := buffer.ReadInt32()
		if err != nil {
			return err
		}
		packet.DimensionID = dimensionID
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
		if err := buffer.WriteInt32(packet.DimensionID); err != nil {
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

	if err := buffer.WriteVarInt(packet.ViewDistance); err != nil {
		return err
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

func (codec DimensionCodec) ToCompound(proto protocol.Protocol) nbt.CompoundTag {
	var dimensions nbt.ListTag
	if proto >= protocol.V1_16_2 {
		for _, dim := range codec.Dimensions {
			dimensions = append(dimensions, nbt.CompoundTag{
				"id":      nbt.IntTag(dim.ID),
				"name":    nbt.StringTag(dim.Name),
				"element": dim.ToCompound(proto),
			})
		}

		var biomes nbt.ListTag
		for _, biome := range codec.Biomes {
			biomes = append(biomes, nbt.CompoundTag{
				"id":      nbt.IntTag(biome.ID),
				"name":    nbt.StringTag(biome.Name),
				"element": biome.ToCompound(),
			})
		}

		return nbt.CompoundTag{
			"minecraft:dimension_type": nbt.CompoundTag{
				"type":  nbt.StringTag("minecraft:dimension_type"),
				"value": dimensions,
			},
			"minecraft:worldgen/biome": nbt.CompoundTag{
				"type":  nbt.StringTag("minecraft:worldgen/biome"),
				"value": biomes,
			},
		}
	}

	for _, dim := range codec.Dimensions {
		dimensions = append(dimensions, dim.ToCompound(proto))
	}
	return nbt.CompoundTag{
		"dimension": dimensions,
	}
}

func (dim Dimension) ToCompound(proto protocol.Protocol) nbt.CompoundTag {
	compound := nbt.CompoundTag{
		"bed_works":            nbt.ByteTag(b2i(dim.BedWorks)),
		"has_ceiling":          nbt.ByteTag(b2i(dim.HasCeiling)),
		"coordinate_scale":     nbt.DoubleTag(dim.CoordinateScale),
		"piglin_safe":          nbt.ByteTag(b2i(dim.PiglinSafe)),
		"has_skylight":         nbt.ByteTag(b2i(dim.HasSkylight)),
		"ultrawarm":            nbt.ByteTag(b2i(dim.Ultrawarm)),
		"infiniburn":           nbt.StringTag(dim.Infiniburn),
		"effects":              nbt.StringTag(dim.Effects),
		"has_raids":            nbt.ByteTag(b2i(dim.HasRaids)),
		"ambient_light":        nbt.FloatTag(dim.AmbientLight),
		"logical_height":       nbt.IntTag(dim.LogicalHeight),
		"natural":              nbt.ByteTag(b2i(dim.Natural)),
		"respawn_anchor_works": nbt.ByteTag(b2i(dim.RespawnAnchorWorks)),
	}

	if dim.FixedTime != nil {
		compound["fixed_time"] = nbt.LongTag(*dim.FixedTime)
	}

	if proto <= protocol.V1_16_1 {
		compound["name"] = nbt.StringTag(dim.Name)
		compound["shrunk"] = nbt.ByteTag(b2i(dim.Shrunk))
	}

	return compound
}

func (biome Biome) ToCompound() nbt.CompoundTag {
	return nbt.CompoundTag{
		"category":      nbt.StringTag(biome.Category),
		"temperature":   nbt.FloatTag(biome.Temperature),
		"downfall":      nbt.FloatTag(biome.Downfall),
		"depth":         nbt.FloatTag(biome.Depth),
		"precipitation": nbt.StringTag(biome.Precipitation),
		"scale":         nbt.FloatTag(biome.Scale),
		"effects": nbt.CompoundTag{
			"water_fog_color": nbt.IntTag(biome.Effects.WaterFogColor),
			"water_color":     nbt.IntTag(biome.Effects.WaterColor),
			"fog_color":       nbt.IntTag(biome.Effects.FogColor),
			"sky_color":       nbt.IntTag(biome.Effects.SkyColor),
			"mood_sound": nbt.CompoundTag{
				"offset":              nbt.DoubleTag(biome.Effects.MoodSound.Offset),
				"sound":               nbt.StringTag(biome.Effects.MoodSound.Sound),
				"block_search_extent": nbt.IntTag(biome.Effects.MoodSound.BlockSearchExtent),
				"tick_delay":          nbt.IntTag(biome.Effects.MoodSound.TickDelay),
			},
		},
	}
}

func DimensionCodecFromTag(tag nbt.Tag, proto protocol.Protocol) (DimensionCodec, error) {
	codec, ok := tag.(nbt.CompoundTag)
	if !ok {
		return DimensionCodec{}, errors.New("tag must be of type nbt.CompoundTag")
	}

	if proto >= protocol.V1_16_2 {
		var dimensions []Dimension
		if value, ok := codec["minecraft:dimension_type"]; ok {
			if value, ok := value.(nbt.CompoundTag); ok {
				if dimList, ok := value["value"]; ok {
					if dimList, ok := dimList.(nbt.ListTag); ok {
						for _, dimInfo := range dimList {
							if dimInfo, ok := dimInfo.(nbt.CompoundTag); ok {
								var id int32
								if value, ok := dimInfo["id"]; ok {
									if tag, ok := value.(nbt.IntTag); ok {
										id = int32(tag)
									}
								}

								var name string
								if value, ok := dimInfo["name"]; ok {
									if tag, ok := value.(nbt.StringTag); ok {
										name = string(tag)
									}
								}

								var dimension Dimension
								if value, ok := dimInfo["element"]; ok {
									dim, err := DimensionFromTag(value)
									if err != nil {
										return DimensionCodec{}, err
									}
									dimension = dim
								}

								dimension.ID = id
								dimension.Name = name
								dimensions = append(dimensions, dimension)
							}
						}
					}
				}
			}
		}

		var biomes []Biome
		if value, ok := codec["minecraft:worldgen/biome"]; ok {
			if value, ok := value.(nbt.CompoundTag); ok {
				if biomeList, ok := value["value"]; ok {
					if biomeList, ok := biomeList.(nbt.ListTag); ok {
						for _, biomeInfo := range biomeList {
							if biomeInfo, ok := biomeInfo.(nbt.CompoundTag); ok {
								var id int32
								if value, ok := biomeInfo["id"]; ok {
									if tag, ok := value.(nbt.IntTag); ok {
										id = int32(tag)
									}
								}

								var name string
								if value, ok := biomeInfo["name"]; ok {
									if tag, ok := value.(nbt.StringTag); ok {
										name = string(tag)
									}
								}

								var biome Biome
								if value, ok := biomeInfo["element"]; ok {
									b, err := BiomeFromTag(value)
									if err != nil {
										return DimensionCodec{}, err
									}
									biome = b
								}

								biome.ID = id
								biome.Name = name
								biomes = append(biomes, biome)
							}
						}
					}
				}
			}
		}

		return DimensionCodec{Dimensions: dimensions, Biomes: biomes}, nil
	}

	var dimensions []Dimension
	if dimList, ok := codec["dimension"]; ok {
		if dimList, ok := dimList.(nbt.ListTag); ok {
			for _, dim := range dimList {
				dimension, err := DimensionFromTag(dim)
				if err != nil {
					return DimensionCodec{}, err
				}
				dimensions = append(dimensions, dimension)
			}
		}
	}
	return DimensionCodec{Dimensions: dimensions}, nil
}

func DimensionFromTag(tag nbt.Tag) (Dimension, error) {
	dim, ok := tag.(nbt.CompoundTag)
	if !ok {
		return Dimension{}, errors.New("tag must be of type nbt.CompoundTag")
	}

	var id int32
	if value, ok := dim["id"]; ok {
		if tag, ok := value.(nbt.IntTag); ok {
			id = int32(tag)
		}
	}

	var name string
	if value, ok := dim["name"]; ok {
		if tag, ok := value.(nbt.StringTag); ok {
			name = string(tag)
		}
	}

	var bedWorks bool
	if value, ok := dim["bed_works"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			bedWorks = i2b(int8(tag))
		}
	}

	var hasCeiling bool
	if value, ok := dim["has_ceiling"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			hasCeiling = i2b(int8(tag))
		}
	}

	var coordinateScale float64
	if value, ok := dim["coordinate_scale"]; ok {
		if tag, ok := value.(nbt.DoubleTag); ok {
			coordinateScale = float64(tag)
		}
	}

	var piglinSafe bool
	if value, ok := dim["piglin_safe"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			piglinSafe = i2b(int8(tag))
		}
	}

	var hasSkylight bool
	if value, ok := dim["has_skylight"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			hasSkylight = i2b(int8(tag))
		}
	}

	var ultrawarm bool
	if value, ok := dim["ultrawarm"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			ultrawarm = i2b(int8(tag))
		}
	}

	var infiniburn string
	if value, ok := dim["infiniburn"]; ok {
		if tag, ok := value.(nbt.StringTag); ok {
			infiniburn = string(tag)
		}
	}

	var effects string
	if value, ok := dim["effects"]; ok {
		if tag, ok := value.(nbt.StringTag); ok {
			effects = string(tag)
		}
	}

	var hasRaids bool
	if value, ok := dim["has_raids"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			hasRaids = i2b(int8(tag))
		}
	}

	var ambientLight float32
	if value, ok := dim["ambient_light"]; ok {
		if tag, ok := value.(nbt.FloatTag); ok {
			ambientLight = float32(tag)
		}
	}

	var logicalHeight int32
	if value, ok := dim["logical_height"]; ok {
		if tag, ok := value.(nbt.IntTag); ok {
			logicalHeight = int32(tag)
		}
	}

	var natural bool
	if value, ok := dim["natural"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			natural = i2b(int8(tag))
		}
	}

	var respawnAnchorWorks bool
	if value, ok := dim["respawn_anchor_works"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			respawnAnchorWorks = i2b(int8(tag))
		}
	}

	var fixedTime *int64
	if value, ok := dim["fixed_time"]; ok {
		if tag, ok := value.(nbt.LongTag); ok {
			tmp := int64(tag)
			fixedTime = &tmp
		}
	}

	var shrunk bool
	if value, ok := dim["shrunk"]; ok {
		if tag, ok := value.(nbt.ByteTag); ok {
			shrunk = i2b(int8(tag))
		}
	}

	return Dimension{
		ID:                 id,
		Name:               name,
		BedWorks:           bedWorks,
		HasCeiling:         hasCeiling,
		CoordinateScale:    coordinateScale,
		PiglinSafe:         piglinSafe,
		HasSkylight:        hasSkylight,
		Ultrawarm:          ultrawarm,
		Infiniburn:         infiniburn,
		Effects:            effects,
		HasRaids:           hasRaids,
		AmbientLight:       ambientLight,
		LogicalHeight:      logicalHeight,
		Natural:            natural,
		RespawnAnchorWorks: respawnAnchorWorks,
		FixedTime:          fixedTime,
		Shrunk:             shrunk,
	}, nil
}

func BiomeFromTag(tag nbt.Tag) (Biome, error) {
	biome, ok := tag.(nbt.CompoundTag)
	if !ok {
		return Biome{}, errors.New("tag must be of type nbt.CompoundTag")
	}

	var id int32
	if value, ok := biome["id"]; ok {
		if tag, ok := value.(nbt.IntTag); ok {
			id = int32(tag)
		}
	}

	var name string
	if value, ok := biome["name"]; ok {
		if tag, ok := value.(nbt.StringTag); ok {
			name = string(tag)
		}
	}

	var category string
	if value, ok := biome["category"]; ok {
		if tag, ok := value.(nbt.StringTag); ok {
			category = string(tag)
		}
	}

	var temperature float32
	if value, ok := biome["temperature"]; ok {
		if tag, ok := value.(nbt.FloatTag); ok {
			temperature = float32(tag)
		}
	}

	var downfall float32
	if value, ok := biome["downfall"]; ok {
		if tag, ok := value.(nbt.FloatTag); ok {
			downfall = float32(tag)
		}
	}

	var depth float32
	if value, ok := biome["depth"]; ok {
		if tag, ok := value.(nbt.FloatTag); ok {
			depth = float32(tag)
		}
	}

	var precipitation string
	if value, ok := biome["precipitation"]; ok {
		if tag, ok := value.(nbt.StringTag); ok {
			precipitation = string(tag)
		}
	}

	var scale float32
	if value, ok := biome["scale"]; ok {
		if tag, ok := value.(nbt.FloatTag); ok {
			scale = float32(tag)
		}
	}

	var effects Effects
	if value, ok := biome["effects"]; ok {
		if tag, ok := value.(nbt.CompoundTag); ok {
			var waterFogColor int32
			if value, ok := tag["water_fog_color"]; ok {
				if tag, ok := value.(nbt.IntTag); ok {
					waterFogColor = int32(tag)
				}
			}

			var waterColor int32
			if value, ok := tag["water_color"]; ok {
				if tag, ok := value.(nbt.IntTag); ok {
					waterColor = int32(tag)
				}
			}

			var fogColor int32
			if value, ok := tag["fog_color"]; ok {
				if tag, ok := value.(nbt.IntTag); ok {
					fogColor = int32(tag)
				}
			}

			var skyColor int32
			if value, ok := tag["sky_color"]; ok {
				if tag, ok := value.(nbt.IntTag); ok {
					skyColor = int32(tag)
				}
			}

			var moodSound MoodSound
			if value, ok := tag["mood_sound"]; ok {
				if tag, ok := value.(nbt.CompoundTag); ok {
					var offset float64
					if value, ok := tag["offset"]; ok {
						if tag, ok := value.(nbt.DoubleTag); ok {
							offset = float64(tag)
						}
					}

					var sound string
					if value, ok := tag["sound"]; ok {
						if tag, ok := value.(nbt.StringTag); ok {
							sound = string(tag)
						}
					}

					var blockSearchExtent int32
					if value, ok := tag["block_search_extent"]; ok {
						if tag, ok := value.(nbt.IntTag); ok {
							blockSearchExtent = int32(tag)
						}
					}

					var tickDelay int32
					if value, ok := tag["tick_delay"]; ok {
						if tag, ok := value.(nbt.IntTag); ok {
							tickDelay = int32(tag)
						}
					}

					moodSound = MoodSound{
						Offset:            offset,
						Sound:             sound,
						BlockSearchExtent: blockSearchExtent,
						TickDelay:         tickDelay,
					}
				}
			}

			effects = Effects{
				WaterFogColor: waterFogColor,
				WaterColor:    waterColor,
				FogColor:      fogColor,
				SkyColor:      skyColor,
				MoodSound:     moodSound,
			}
		}
	}

	return Biome{
		ID:            id,
		Name:          name,
		Category:      category,
		Temperature:   temperature,
		Downfall:      downfall,
		Depth:         depth,
		Precipitation: precipitation,
		Scale:         scale,
		Effects:       effects,
	}, nil
}

func b2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func i2b(i int8) bool {
	if i == 1 {
		return true
	}
	return false
}
