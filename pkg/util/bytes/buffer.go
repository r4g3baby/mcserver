package bytes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	"io"
	"math"
)

type Buffer struct {
	*bytes.Buffer
}

func (buffer *Buffer) ReadVarInt() (int32, error) {
	var result int32 = 0
	for numRead := 0; ; numRead++ {
		read, err := buffer.ReadByte()
		if err != nil {
			return result, err
		}

		result |= int32(read&0x7F) << (7 * numRead)

		if numRead >= 5 {
			return result, errors.New("VarInt too big")
		}

		if (read & 0x80) != 0x80 {
			break
		}
	}
	return result, nil
}

func (buffer *Buffer) ReadVarLong() (int64, error) {
	var result int64 = 0
	for numRead := 0; ; numRead++ {
		read, err := buffer.ReadByte()
		if err != nil {
			return result, err
		}

		result |= int64(read&0x7F) << (7 * numRead)

		if numRead >= 10 {
			return result, errors.New("VarLong too big")
		}

		if (read & 0x80) != 0x80 {
			break
		}
	}
	return result, nil
}

func (buffer *Buffer) ReadUtf(maxLength int) (string, error) {
	length, err := buffer.ReadVarInt()
	if err != nil {
		return "", err
	}

	if length < 1 || int(length) > (maxLength*4)+3 {
		return "", errors.New("the received encoded string bytes length is invalid")
	}

	var str = make([]byte, length)
	n, err := io.ReadFull(buffer, str)
	if err != nil {
		return "", err
	}

	if n > maxLength {
		return "", errors.New("the received string length is longer than maximum allowed")
	}

	return string(str), nil
}

func (buffer *Buffer) ReadUUID() (uuid.UUID, error) {
	var uuidBytes = make([]byte, 16)
	_, err := io.ReadFull(buffer, uuidBytes)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(uuidBytes)
}

func (buffer *Buffer) ReadBool() (bool, error) {
	value, err := buffer.ReadUint8()
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

func (buffer *Buffer) ReadInt8() (int8, error) {
	value, err := buffer.ReadUint8()
	if err != nil {
		return 0, err
	}
	return int8(value), nil
}

func (buffer *Buffer) ReadUint8() (uint8, error) {
	var value = make([]byte, 1)
	_, err := buffer.Read(value)
	if err != nil {
		return 0, err
	}
	return value[0], nil
}

func (buffer *Buffer) ReadInt16() (int16, error) {
	value, err := buffer.ReadUint16()
	if err != nil {
		return 0, err
	}
	return int16(value), nil
}

func (buffer *Buffer) ReadUint16() (uint16, error) {
	var value = make([]byte, 2)
	_, err := buffer.Read(value)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(value), nil
}

func (buffer *Buffer) ReadInt32() (int32, error) {
	value, err := buffer.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

func (buffer *Buffer) ReadUint32() (uint32, error) {
	var value = make([]byte, 4)
	_, err := buffer.Read(value)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(value), nil
}

func (buffer *Buffer) ReadInt64() (int64, error) {
	value, err := buffer.ReadUint64()
	if err != nil {
		return 0, err
	}
	return int64(value), nil
}

func (buffer *Buffer) ReadUint64() (uint64, error) {
	var value = make([]byte, 8)
	_, err := buffer.Read(value)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(value), nil
}

func (buffer *Buffer) ReadFloat32() (float32, error) {
	value, err := buffer.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(value), nil
}

func (buffer *Buffer) ReadFloat64() (float64, error) {
	value, err := buffer.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(value), nil
}

func (buffer *Buffer) WriteVarInt(value int32) error {
	for value >= 0x80 {
		err := buffer.WriteByte(byte(value) | 0x80)
		if err != nil {
			return err
		}
		value >>= 7
	}
	return buffer.WriteByte(byte(value))
}

func (buffer *Buffer) WriteVarLong(value int64) error {
	for value >= 0x80 {
		err := buffer.WriteByte(byte(value) | 0x80)
		if err != nil {
			return err
		}
		value >>= 7
	}
	return buffer.WriteByte(byte(value))
}

func (buffer *Buffer) WriteUtf(value string, maxLength int) error {
	valueB := []byte(value)
	if len(valueB) > maxLength {
		return errors.New("string too big")
	}
	err := buffer.WriteVarInt(int32(len(valueB)))
	if err != nil {
		return err
	}
	_, err = buffer.Write(valueB)
	return err
}

func (buffer *Buffer) WriteUUID(uuid uuid.UUID) error {
	err := buffer.WriteUint64(binary.BigEndian.Uint64(uuid[:8]))
	if err != nil {
		return err
	}
	return buffer.WriteUint64(binary.BigEndian.Uint64(uuid[8:]))
}

func (buffer *Buffer) WriteBool(value bool) error {
	if value {
		return buffer.WriteUint8(1)
	} else {
		return buffer.WriteUint8(0)
	}
}

func (buffer *Buffer) WriteInt8(value int8) error {
	return buffer.WriteUint8(uint8(value))
}

func (buffer *Buffer) WriteUint8(value uint8) error {
	var valueB = make([]byte, 1)
	valueB[0] = value
	_, err := buffer.Write(valueB)
	return err
}

func (buffer *Buffer) WriteInt16(value int16) error {
	return buffer.WriteUint16(uint16(value))
}

func (buffer *Buffer) WriteUint16(value uint16) error {
	var valueB = make([]byte, 2)
	binary.BigEndian.PutUint16(valueB, value)
	_, err := buffer.Write(valueB)
	return err
}

func (buffer *Buffer) WriteInt32(value int32) error {
	return buffer.WriteUint32(uint32(value))
}

func (buffer *Buffer) WriteUint32(value uint32) error {
	var valueB = make([]byte, 4)
	binary.BigEndian.PutUint32(valueB, value)
	_, err := buffer.Write(valueB)
	return err
}

func (buffer *Buffer) WriteInt64(value int64) error {
	return buffer.WriteUint64(uint64(value))
}

func (buffer *Buffer) WriteUint64(value uint64) error {
	var valueB = make([]byte, 8)
	binary.BigEndian.PutUint64(valueB, value)
	_, err := buffer.Write(valueB)
	return err
}

func (buffer *Buffer) WriteFloat32(value float32) error {
	var valueB = make([]byte, 4)
	binary.BigEndian.PutUint32(valueB, math.Float32bits(value))
	_, err := buffer.Write(valueB)
	return err
}

func (buffer *Buffer) WriteFloat64(value float64) error {
	var valueB = make([]byte, 8)
	binary.BigEndian.PutUint64(valueB, math.Float64bits(value))
	_, err := buffer.Write(valueB)
	return err
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{bytes.NewBuffer(buf)}
}
