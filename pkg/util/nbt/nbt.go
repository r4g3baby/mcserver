package nbt

type Type uint8

const (
	TypeEnd Type = iota
	TypeByte
	TypeShort
	TypeInt
	TypeLong
	TypeFloat
	TypeDouble
	TypeByteArray
	TypeString
	TypeList
	TypeCompound
	TypeIntArray
	TypeLongArray
)

type (
	Tag interface {
		Type() Type
	}

	EndTag       struct{}
	ByteTag      uint8
	ShortTag     int16
	IntTag       int32
	LongTag      int64
	FloatTag     float32
	DoubleTag    float64
	ByteArrayTag []uint8
	StringTag    string
	ListTag      []Tag
	CompoundTag  map[string]Tag
	IntArrayTag  []int32
	LongArrayTag []int64
)

func (tag EndTag) Type() Type {
	return TypeEnd
}

func (tag ByteTag) Type() Type {
	return TypeByte
}

func (tag ShortTag) Type() Type {
	return TypeShort
}

func (tag IntTag) Type() Type {
	return TypeInt
}

func (tag LongTag) Type() Type {
	return TypeLong
}

func (tag FloatTag) Type() Type {
	return TypeFloat
}

func (tag DoubleTag) Type() Type {
	return TypeDouble
}

func (tag ByteArrayTag) Type() Type {
	return TypeByteArray
}

func (tag StringTag) Type() Type {
	return TypeString
}

func (tag ListTag) Type() Type {
	return TypeList
}

func (tag CompoundTag) Type() Type {
	return TypeCompound
}

func (tag IntArrayTag) Type() Type {
	return TypeIntArray
}

func (tag LongArrayTag) Type() Type {
	return TypeLongArray
}
