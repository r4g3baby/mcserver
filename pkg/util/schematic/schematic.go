package schematic

type (
	Schematic interface {
		GetVersion() int
		GetDataVersion() int
		GetMetadata() Metadata
		GetWidth() int
		GetHeight() int
		GetLength() int
	}

	Metadata interface {
		GetName() string
		GetAuthor() string
		GetDate() int64
		GetRequiredMods() []string
	}
)

type (
	schematic struct {
		version, dataVersion  int
		metadata              Metadata
		width, height, length int
		offset                [3]int
		blocks                [][][]string
	}

	metadata struct {
		name, author string
		date         int64
		requiredMods []string
	}
)

func (schem *schematic) GetVersion() int {
	return schem.version
}

func (schem *schematic) GetDataVersion() int {
	return schem.dataVersion
}

func (schem *schematic) GetMetadata() Metadata {
	return schem.metadata
}

func (schem *schematic) GetWidth() int {
	return schem.width
}

func (schem *schematic) GetHeight() int {
	return schem.height
}

func (schem *schematic) GetLength() int {
	return schem.length
}

func (meta *metadata) GetName() string {
	return meta.name
}

func (meta *metadata) GetAuthor() string {
	return meta.author
}

func (meta *metadata) GetDate() int64 {
	return meta.date
}

func (meta *metadata) GetRequiredMods() []string {
	return meta.requiredMods
}
