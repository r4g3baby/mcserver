package nbt

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func ReadCompressed(reader io.Reader) (string, Tag, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", nil, err
	}
	defer gzipReader.Close()
	return Read(gzipReader)
}

func Read(reader io.Reader) (string, Tag, error) {
	typeByte, err := readByte(reader)
	if err != nil {
		return "", nil, err
	}

	typ := Type(typeByte)
	if typ == TypeEnd {
		return "", EndTag{}, nil
	}

	name, err := readString(reader)
	if err != nil {
		return "", nil, err
	}

	switch typ {
	case TypeByte:
		tag, err := readByte(reader)
		return string(name), tag, err
	case TypeShort:
		tag, err := readShort(reader)
		return string(name), tag, err
	case TypeInt:
		tag, err := readInt(reader)
		return string(name), tag, err
	case TypeLong:
		tag, err := readLong(reader)
		return string(name), tag, err
	case TypeFloat:
		tag, err := readFloat(reader)
		return string(name), tag, err
	case TypeDouble:
		tag, err := readDouble(reader)
		return string(name), tag, err
	case TypeByteArray:
		tag, err := readByteArray(reader)
		return string(name), tag, err
	case TypeString:
		tag, err := readString(reader)
		return string(name), tag, err
	case TypeList:
		tag, err := readList(reader)
		return string(name), tag, err
	case TypeCompound:
		tag, err := readCompound(reader)
		return string(name), tag, err
	case TypeIntArray:
		tag, err := readIntArray(reader)
		return string(name), tag, err
	case TypeLongArray:
		tag, err := readLongArray(reader)
		return string(name), tag, err
	default:
		return string(name), nil, fmt.Errorf("unsupported tag type %v", typ)
	}
}

func readByte(reader io.Reader) (ByteTag, error) {
	var buff = make([]byte, 1)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}
	return ByteTag(buff[0]), nil
}

func readShort(reader io.Reader) (ShortTag, error) {
	var buff = make([]byte, 2)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}
	return ShortTag(binary.BigEndian.Uint16(buff)), nil
}

func readInt(reader io.Reader) (IntTag, error) {
	var buff = make([]byte, 4)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}
	return IntTag(binary.BigEndian.Uint32(buff)), nil
}

func readLong(reader io.Reader) (LongTag, error) {
	var buff = make([]byte, 8)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}
	return LongTag(binary.BigEndian.Uint64(buff)), nil
}

func readFloat(reader io.Reader) (FloatTag, error) {
	var buff = make([]byte, 4)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}
	return FloatTag(math.Float32frombits(binary.BigEndian.Uint32(buff))), nil
}

func readDouble(reader io.Reader) (DoubleTag, error) {
	var buff = make([]byte, 8)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}
	return DoubleTag(math.Float64frombits(binary.BigEndian.Uint64(buff))), nil
}

func readByteArray(reader io.Reader) (ByteArrayTag, error) {
	size, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	var byteArray ByteArrayTag
	for i := size; i > 0; i-- {
		b, err := readByte(reader)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, byte(b))
	}
	return byteArray, nil
}

func readString(reader io.Reader) (StringTag, error) {
	length, err := readShort(reader)
	if err != nil {
		return "", err
	}

	var buff = make([]byte, length)
	if _, err = reader.Read(buff); err != nil {
		return "", err
	}
	return StringTag(buff), nil
}

func readList(reader io.Reader) (ListTag, error) {
	typeByte, err := readByte(reader)
	if err != nil {
		return nil, err
	}
	typ := Type(typeByte)

	count, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	var list ListTag
	switch typ {
	case TypeEnd:
		break
	case TypeByte:
		for i := count; i > 0; i-- {
			tag, err := readByte(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeShort:
		for i := count; i > 0; i-- {
			tag, err := readShort(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeInt:
		for i := count; i > 0; i-- {
			tag, err := readInt(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeLong:
		for i := count; i > 0; i-- {
			tag, err := readLong(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeFloat:
		for i := count; i > 0; i-- {
			tag, err := readFloat(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeDouble:
		for i := count; i > 0; i-- {
			tag, err := readDouble(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeByteArray:
		for i := count; i > 0; i-- {
			tag, err := readByteArray(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeString:
		for i := count; i > 0; i-- {
			tag, err := readString(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeList:
		for i := count; i > 0; i-- {
			tag, err := readList(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeCompound:
		for i := count; i > 0; i-- {
			tag, err := readCompound(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeIntArray:
		for i := count; i > 0; i-- {
			tag, err := readIntArray(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	case TypeLongArray:
		for i := count; i > 0; i-- {
			tag, err := readLongArray(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, tag)
		}
	default:
		return nil, fmt.Errorf("unsupported tag type %v", typ)
	}
	return list, nil
}

func readCompound(reader io.Reader) (CompoundTag, error) {
	var compound = make(CompoundTag)
	for name, tag, err := Read(reader); tag != nil && tag.Type() != TypeEnd; name, tag, err = Read(reader) {
		if err != nil {
			return nil, err
		}
		compound[name] = tag
	}
	return compound, nil
}

func readIntArray(reader io.Reader) (IntArrayTag, error) {
	size, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	var intArray IntArrayTag
	for i := size; i > 0; i-- {
		b, err := readInt(reader)
		if err != nil {
			return nil, err
		}
		intArray = append(intArray, int32(b))
	}
	return intArray, nil
}

func readLongArray(reader io.Reader) (LongArrayTag, error) {
	size, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	var longArray LongArrayTag
	for i := size; i > 0; i-- {
		b, err := readLong(reader)
		if err != nil {
			return nil, err
		}
		longArray = append(longArray, int64(b))
	}
	return longArray, nil
}
