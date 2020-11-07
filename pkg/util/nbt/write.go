package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func Write(writer io.Writer, name string, tag Tag) error {
	if err := writeByte(writer, ByteTag(tag.Type())); err != nil {
		return err
	}

	if err := writeString(writer, StringTag(name)); err != nil {
		return err
	}

	switch t := tag.(type) {
	case EndTag:
		return nil
	case ByteTag:
		return writeByte(writer, t)
	case ShortTag:
		return writeShort(writer, t)
	case IntTag:
		return writeInt(writer, t)
	case LongTag:
		return writeLong(writer, t)
	case FloatTag:
		return writeFloat(writer, t)
	case DoubleTag:
		return writeDouble(writer, t)
	case ByteArrayTag:
		return writeByteArray(writer, t)
	case StringTag:
		return writeString(writer, t)
	case ListTag:
		return writeList(writer, t)
	case CompoundTag:
		return writeCompound(writer, t)
	case IntArrayTag:
		return writeIntArray(writer, t)
	case LongArrayTag:
		return writeLongArray(writer, t)
	default:
		return fmt.Errorf("unsupported tag %T", tag)
	}
}

func writeByte(writer io.Writer, value ByteTag) error {
	var buff = make([]byte, 1)
	buff[0] = byte(value)
	_, err := writer.Write(buff)
	return err
}

func writeShort(writer io.Writer, value ShortTag) error {
	var buff = make([]byte, 2)
	binary.BigEndian.PutUint16(buff, uint16(value))
	_, err := writer.Write(buff)
	return err
}

func writeInt(writer io.Writer, value IntTag) error {
	var buff = make([]byte, 4)
	binary.BigEndian.PutUint32(buff, uint32(value))
	_, err := writer.Write(buff)
	return err
}

func writeLong(writer io.Writer, value LongTag) error {
	var buff = make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(value))
	_, err := writer.Write(buff)
	return err
}

func writeFloat(writer io.Writer, value FloatTag) error {
	var buff = make([]byte, 4)
	binary.BigEndian.PutUint32(buff, math.Float32bits(float32(value)))
	_, err := writer.Write(buff)
	return err
}

func writeDouble(writer io.Writer, value DoubleTag) error {
	var buff = make([]byte, 8)
	binary.BigEndian.PutUint64(buff, math.Float64bits(float64(value)))
	_, err := writer.Write(buff)
	return err
}

func writeByteArray(writer io.Writer, value ByteArrayTag) error {
	if err := writeInt(writer, IntTag(len(value))); err != nil {
		return err
	}
	for _, b := range value {
		if err := writeByte(writer, ByteTag(b)); err != nil {
			return err
		}
	}
	return nil
}

func writeString(writer io.Writer, value StringTag) error {
	buff := []byte(value)
	if err := writeShort(writer, ShortTag(len(buff))); err != nil {
		return err
	}
	_, err := writer.Write(buff)
	return err
}

func writeList(writer io.Writer, value ListTag) error {
	if len(value) == 0 {
		if err := writeByte(writer, ByteTag(TypeEnd)); err != nil {
			return err
		}
		return writeInt(writer, IntTag(0))
	}

	typ := value[0].Type()
	if err := writeByte(writer, ByteTag(typ)); err != nil {
		return err
	}
	if err := writeInt(writer, IntTag(len(value))); err != nil {
		return err
	}

	switch typ {
	case TypeEnd:
		break
	case TypeByte:
		for _, t := range value {
			if tag, ok := t.(ByteTag); ok {
				if err := writeByte(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeShort:
		for _, t := range value {
			if tag, ok := t.(ShortTag); ok {
				if err := writeShort(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeInt:
		for _, t := range value {
			if tag, ok := t.(IntTag); ok {
				if err := writeInt(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeLong:
		for _, t := range value {
			if tag, ok := t.(LongTag); ok {
				if err := writeLong(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeFloat:
		for _, t := range value {
			if tag, ok := t.(FloatTag); ok {
				if err := writeFloat(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeDouble:
		for _, t := range value {
			if tag, ok := t.(DoubleTag); ok {
				if err := writeDouble(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeByteArray:
		for _, t := range value {
			if tag, ok := t.(ByteArrayTag); ok {
				if err := writeByteArray(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeString:
		for _, t := range value {
			if tag, ok := t.(StringTag); ok {
				if err := writeString(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeList:
		for _, t := range value {
			if tag, ok := t.(ListTag); ok {
				if err := writeList(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeCompound:
		for _, t := range value {
			if tag, ok := t.(CompoundTag); ok {
				if err := writeCompound(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeIntArray:
		for _, t := range value {
			if tag, ok := t.(IntArrayTag); ok {
				if err := writeIntArray(writer, tag); err != nil {
					return err
				}
			}
		}
	case TypeLongArray:
		for _, t := range value {
			if tag, ok := t.(LongArrayTag); ok {
				if err := writeLongArray(writer, tag); err != nil {
					return err
				}
			}
		}
	default:
		return fmt.Errorf("unsupported tag type %v", typ)
	}
	return nil
}

func writeCompound(writer io.Writer, value CompoundTag) error {
	for name, tag := range value {
		if err := Write(writer, name, tag); err != nil {
			return err
		}
	}
	return writeByte(writer, ByteTag(TypeEnd))
}

func writeIntArray(writer io.Writer, value IntArrayTag) error {
	if err := writeInt(writer, IntTag(len(value))); err != nil {
		return err
	}
	for _, b := range value {
		if err := writeInt(writer, IntTag(b)); err != nil {
			return err
		}
	}
	return nil
}

func writeLongArray(writer io.Writer, value LongArrayTag) error {
	if err := writeInt(writer, IntTag(len(value))); err != nil {
		return err
	}
	for _, b := range value {
		if err := writeLong(writer, LongTag(b)); err != nil {
			return err
		}
	}
	return nil
}
