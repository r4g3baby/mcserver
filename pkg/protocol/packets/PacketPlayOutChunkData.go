package packets

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/nbt"
)

type PacketPlayOutChunkData struct {
	ChunkX, ChunkZ int32
	FullChunk      bool
	PrimaryBit     int32
	Heightmaps     nbt.Tag
	Biomes         []int32
	Data           []byte
	BlockEntities  []nbt.Tag
}

func (packet *PacketPlayOutChunkData) GetID(proto protocol.Protocol) (int32, error) {
	return GetID(proto, protocol.Play, protocol.ClientBound, packet)
}

func (packet *PacketPlayOutChunkData) Read(_ protocol.Protocol, buffer *bytes.Buffer) error {
	chunkX, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	packet.ChunkX = chunkX

	chunkZ, err := buffer.ReadInt32()
	if err != nil {
		return err
	}
	packet.ChunkZ = chunkZ

	fullChunk, err := buffer.ReadBool()
	if err != nil {
		return err
	}
	packet.FullChunk = fullChunk

	primaryBit, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}
	packet.PrimaryBit = primaryBit

	_, heightmaps, err := nbt.Read(buffer)
	if err != nil {
		return err
	}
	packet.Heightmaps = heightmaps

	if packet.FullChunk {
		biomesCount, err := buffer.ReadVarInt()
		if err != nil {
			return err
		}

		var biomes []int32
		for i := biomesCount; i > 0; i-- {
			biome, err := buffer.ReadVarInt()
			if err != nil {
				return err
			}
			biomes = append(biomes, biome)
		}
		packet.Biomes = biomes
	}

	size, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}

	var data = make([]byte, size)
	_, err = buffer.Read(data)
	if err != nil {
		return err
	}
	packet.Data = data

	blockEntitiesCount, err := buffer.ReadVarInt()
	if err != nil {
		return err
	}

	var blockEntities []nbt.Tag
	for i := blockEntitiesCount; i > 0; i-- {
		_, blockEntity, err := nbt.Read(buffer)
		if err != nil {
			return err
		}
		blockEntities = append(blockEntities, blockEntity)
	}
	packet.BlockEntities = blockEntities

	return nil
}

func (packet *PacketPlayOutChunkData) Write(_ protocol.Protocol, buffer *bytes.Buffer) error {
	if err := buffer.WriteInt32(packet.ChunkX); err != nil {
		return err
	}

	if err := buffer.WriteInt32(packet.ChunkZ); err != nil {
		return err
	}

	if err := buffer.WriteBool(packet.FullChunk); err != nil {
		return err
	}

	if err := buffer.WriteVarInt(packet.PrimaryBit); err != nil {
		return err
	}

	if err := nbt.Write(buffer, "", packet.Heightmaps); err != nil {
		return err
	}

	if packet.FullChunk {
		if err := buffer.WriteVarInt(int32(len(packet.Biomes))); err != nil {
			return err
		}

		for _, biome := range packet.Biomes {
			if err := buffer.WriteVarInt(biome); err != nil {
				return err
			}
		}
	}

	if err := buffer.WriteVarInt(int32(len(packet.Data))); err != nil {
		return err
	}

	if _, err := buffer.Write(packet.Data); err != nil {
		return err
	}

	if err := buffer.WriteVarInt(int32(len(packet.BlockEntities))); err != nil {
		return err
	}

	for _, blockEntity := range packet.BlockEntities {
		if err := nbt.Write(buffer, "", blockEntity); err != nil {
			return err
		}
	}

	return nil
}
