package schematic

import (
	"errors"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/nbt"
	"io"
)

func Read(reader io.Reader) (Schematic, error) {
	_, root, err := nbt.ReadCompressed(reader)
	if err != nil {
		return nil, err
	}

	tag, ok := root.(nbt.CompoundTag)
	if !ok {
		return nil, errors.New("root tag must be of type nbt.CompoundTag")
	}

	var schem schematic
	schem.version = int(tag["Version"].(nbt.IntTag))
	schem.dataVersion = int(tag["DataVersion"].(nbt.IntTag))

	if meta, ok := tag["Metadata"]; ok {
		meta := meta.(nbt.CompoundTag)

		var name string
		if mName, ok := meta["Name"]; ok {
			name = string(mName.(nbt.StringTag))
		}

		var author string
		if mAuthor, ok := meta["Author"]; ok {
			author = string(mAuthor.(nbt.StringTag))
		}

		var date int64
		if mDate, ok := meta["Date"]; ok {
			date = int64(mDate.(nbt.LongTag))
		}

		schem.metadata = &metadata{
			name:         name,
			author:       author,
			date:         date,
			requiredMods: nil,
		}
	}

	schem.width = int(tag["Width"].(nbt.ShortTag))
	schem.height = int(tag["Height"].(nbt.ShortTag))
	schem.length = int(tag["Length"].(nbt.ShortTag))

	schem.blocks = make([][][]string, schem.width)
	for x := range schem.blocks {
		schem.blocks[x] = make([][]string, schem.height)
		for y := range schem.blocks[x] {
			schem.blocks[x][y] = make([]string, schem.length)
		}
	}

	paletteMax := tag["PaletteMax"].(nbt.IntTag)
	paletteObj := tag["Palette"].(nbt.CompoundTag)
	if int(paletteMax) != len(paletteObj) {
		return nil, errors.New("block palette size does not match expected size")
	}

	var palette = make(map[int32]string)
	for block, index := range paletteObj {
		palette[int32(index.(nbt.IntTag))] = block
	}

	blocks := tag["BlockData"].(nbt.ByteArrayTag)
	buff := bytes.NewBuffer(blocks)
	for i := 0; i < len(blocks); i++ {
		paletteIndex, err := buff.ReadVarInt()
		if err != nil {
			panic(err)
		}

		y := i / (schem.width * schem.length)
		x := (i % (schem.width * schem.length)) % schem.width
		z := (i % (schem.width * schem.length)) / schem.width
		schem.blocks[x][y][z] = palette[paletteIndex]
	}

	return &schem, nil
}
