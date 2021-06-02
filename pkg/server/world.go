package server

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/blocks"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/nbt"
	"github.com/r4g3baby/mcserver/pkg/util/pools"
)

const (
	ChunkWidth    = 16
	ChunkHeight   = 256
	SectionWidth  = ChunkWidth
	SectionHeight = 16
	ChunkSections = ChunkHeight / SectionHeight
	SectionVolume = SectionWidth * SectionHeight * SectionWidth

	MinBitsPerBlock    = 4
	MaxBitsPerBlock    = 8
	GlobalBitsPerBlock = 14
)

type (
	World interface {
		GetName() string
		GetDimension() protocol.Dimension
		GetChunk(x, z int) Chunk
		GetChunks() []Chunk
		SetBlock(x, y, z int, block string)
		GetBlock(x, y, z int) string
		SendChunks(player Player) error
	}

	Chunk interface {
		GetX() int
		GetZ() int
		GetSection(y int) ChunkSection
		GetSections() [ChunkSections]ChunkSection
		SetBlock(x, y, z int, block string)
		GetBlock(x, y, z int) string
	}

	ChunkSection interface {
		IsEmpty() bool
		GetPalette() SectionPalette
		SetBlock(x, y, z int, block string)
		GetBlock(x, y, z int) string
		GetBlocks() bytes.PackedArray
	}

	SectionPalette interface {
		GetOrAdd(block string) int
		GetIndex(block string) int
		GetBlock(index int) string
		GetBlocks() []string
		GetLength() int
		GetBitsPerBlock() int
	}
)

type (
	world struct {
		name      string
		dimension protocol.Dimension
		chunks    []Chunk
	}

	chunk struct {
		x, z     int
		sections [ChunkSections]ChunkSection
	}

	chunkSection struct {
		palette SectionPalette
		blocks  bytes.PackedArray
	}

	sectionPalette struct {
		blocks []string
	}
)

func (world *world) GetName() string {
	return world.name
}

func (world *world) GetDimension() protocol.Dimension {
	return world.dimension
}

func (world *world) GetChunk(x, z int) Chunk {
	for _, chunk := range world.chunks {
		if chunk.GetX() == x && chunk.GetZ() == z {
			return chunk
		}
	}

	chunk := &chunk{x: x, z: z}
	world.chunks = append(world.chunks, chunk)
	return chunk
}

func (world *world) GetChunks() []Chunk {
	return world.chunks
}

func (world *world) SetBlock(x, y, z int, block string) {
	world.GetChunk(x>>4, z>>4).SetBlock(mod(x, 16), y, mod(z, 16), block)
}

func (world *world) GetBlock(x, y, z int) string {
	return world.GetChunk(x>>4, z>>4).GetBlock(mod(x, 16), y, mod(z, 16))
}

func (world *world) SendChunks(player Player) error {
	var biomes []int32
	for i := 0; i < 1024; i++ {
		biomes = append(biomes, 127)
	}

	for _, chunk := range world.chunks {
		if err := func(data *bytes.Buffer) error {
			defer pools.Buffer.Put(data)

			mask := 0
			for sectionY, section := range chunk.GetSections() {
				if section != nil && !section.IsEmpty() {
					mask |= 1 << sectionY

					if err := data.WriteUint16(SectionVolume); err != nil {
						return err
					}

					bitsPerBlock := section.GetBlocks().GetBitsPerValue()
					if err := data.WriteUint8(uint8(bitsPerBlock)); err != nil {
						return err
					}

					var blockData []uint64
					if bitsPerBlock == GlobalBitsPerBlock {
						// Since we depend on the player protocol version we have to create the blockData here
						// This is a bit more expensive but it's the best approach I can think of atm
						blocksArray := bytes.NewPackedArray(bitsPerBlock, SectionVolume)
						for i := 0; i < section.GetBlocks().GetCapacity(); i++ {
							block := section.GetPalette().GetBlock(section.GetBlocks().Get(i))
							blocksArray.Set(i, blocks.GetBlockID(block, player.GetProtocol()))
						}

						blockData = blocksArray.GetData()
					} else {
						if err := data.WriteVarInt(int32(section.GetPalette().GetLength())); err != nil {
							return err
						}
						for _, block := range section.GetPalette().GetBlocks() {
							id := blocks.GetBlockID(block, player.GetProtocol())
							if err := data.WriteVarInt(int32(id)); err != nil {
								return err
							}
						}

						blockData = section.GetBlocks().GetData()
					}

					if err := data.WriteVarInt(int32(len(blockData))); err != nil {
						return err
					}
					for _, value := range blockData {
						if err := data.WriteUint64(value); err != nil {
							return err
						}
					}
				}
			}

			return player.SendPacket(&packets.PacketPlayOutChunkData{
				ChunkX:        int32(chunk.GetX()),
				ChunkZ:        int32(chunk.GetZ()),
				FullChunk:     true,
				PrimaryBit:    int32(mask),
				Heightmaps:    nbt.CompoundTag{},
				Biomes:        biomes,
				Data:          data.Bytes(),
				BlockEntities: []nbt.Tag{},
			})
		}(pools.Buffer.Get(nil)); err != nil {
			return err
		}
	}

	return nil
}

func (chunk *chunk) GetX() int {
	return chunk.x
}

func (chunk *chunk) GetZ() int {
	return chunk.z
}

func (chunk *chunk) GetSections() [ChunkSections]ChunkSection {
	return chunk.sections
}

func (chunk *chunk) GetSection(y int) ChunkSection {
	if y < 0 || y > ChunkSections {
		return nil
	}

	section := chunk.sections[y]
	if section != nil {
		return section
	}

	section = &chunkSection{
		palette: &sectionPalette{[]string{"minecraft:air"}},
		blocks:  bytes.NewPackedArray(MinBitsPerBlock, SectionVolume),
	}
	chunk.sections[y] = section
	return section
}

func (chunk *chunk) SetBlock(x, y, z int, block string) {
	section := chunk.GetSection(y >> 4)
	if section == nil {
		return
	}

	section.SetBlock(x, mod(y, 16), z, block)
}

func (chunk *chunk) GetBlock(x, y, z int) string {
	section := chunk.GetSection(y >> 4)
	if section == nil {
		return "minecraft:air"
	}

	return section.GetBlock(x, mod(y, 16), z)
}

func (section *chunkSection) IsEmpty() bool {
	return section.palette.GetLength() == 1
}

func (section *chunkSection) GetPalette() SectionPalette {
	return section.palette
}

func (section *chunkSection) SetBlock(x, y, z int, block string) {
	id := section.palette.GetOrAdd(block)
	if section.blocks.GetBitsPerValue() != section.palette.GetBitsPerBlock() {
		section.blocks = section.blocks.Resized(section.palette.GetBitsPerBlock())
	}
	section.blocks.Set(index(x, y, z), id)
}

func (section *chunkSection) GetBlock(x, y, z int) string {
	return section.palette.GetBlock(section.blocks.Get(index(x, y, z)))
}

func (section *chunkSection) GetBlocks() bytes.PackedArray {
	return section.blocks
}

func (sPalette *sectionPalette) GetOrAdd(block string) int {
	for i, v := range sPalette.blocks {
		if v == block {
			return i
		}
	}
	sPalette.blocks = append(sPalette.blocks, block)
	return len(sPalette.blocks) - 1
}

func (sPalette *sectionPalette) GetIndex(block string) int {
	for i, b := range sPalette.blocks {
		if b == block {
			return i
		}
	}
	return 0
}

func (sPalette *sectionPalette) GetBlock(id int) string {
	for i, b := range sPalette.blocks {
		if i == id {
			return b
		}
	}
	return "minecraft:air"
}

func (sPalette *sectionPalette) GetBlocks() []string {
	return sPalette.blocks
}

func (sPalette *sectionPalette) GetLength() int {
	return len(sPalette.blocks)
}

func (sPalette *sectionPalette) GetBitsPerBlock() int {
	bitsPerBlock := 0
	for n := len(sPalette.blocks); n != 0; n >>= 1 {
		bitsPerBlock++
	}

	if bitsPerBlock < MinBitsPerBlock {
		bitsPerBlock = MinBitsPerBlock
	} else if bitsPerBlock > MaxBitsPerBlock {
		bitsPerBlock = GlobalBitsPerBlock
	}
	return bitsPerBlock
}

func NewWorld(name string, dimension protocol.Dimension) World {
	return &world{
		name:      name,
		dimension: dimension,
	}
}

func index(x, y, z int) int {
	return (y&0xf)<<8 | z<<4 | x
}

func mod(a, b int) int {
	return (a%b + b) % b
}
