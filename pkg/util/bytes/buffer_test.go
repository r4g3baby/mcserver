package bytes

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

var buffer = NewBuffer(nil)

func cleanup() {
	buffer.Truncate(0)
}

func TestBuffer_VarInt(t *testing.T) {
	t.Cleanup(cleanup)
	var want int32 = 1337

	if err := buffer.WriteVarInt(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadVarInt()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("VarInt was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_VarLong(t *testing.T) {
	t.Cleanup(cleanup)
	var want int64 = 1337

	if err := buffer.WriteVarLong(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadVarLong()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("VarLong was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Utf(t *testing.T) {
	t.Cleanup(cleanup)
	var want = "Cats don't always land on their feet."

	if err := buffer.WriteUtf(want, 37); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadUtf(37)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Utf was incorrect, got: %s, want: %s.", got, want)
	}
}

func TestBuffer_UUID(t *testing.T) {
	t.Cleanup(cleanup)
	var want, err = uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}

	if err := buffer.WriteUUID(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadUUID()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("UUID was incorrect, got: %s, want: %s.", got, want)
	}
}

func TestBuffer_Bool(t *testing.T) {
	t.Cleanup(cleanup)
	for _, want := range []bool{true, false} {
		t.Run(fmt.Sprintf("%v", want), func(t *testing.T) {
			if err := buffer.WriteBool(want); err != nil {
				t.Fatal(err)
			}

			got, err := buffer.ReadBool()
			if err != nil {
				t.Fatal(err)
			}

			if got != want {
				t.Errorf("Bool was incorrect, got: %v, want: %v.", got, want)
			}
		})
	}
}

func TestBuffer_Int8(t *testing.T) {
	t.Cleanup(cleanup)
	var want int8 = 23

	if err := buffer.WriteInt8(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadInt8()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Int8 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Uint8(t *testing.T) {
	t.Cleanup(cleanup)
	var want uint8 = 23

	if err := buffer.WriteUint8(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadUint8()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Uint8 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Int16(t *testing.T) {
	t.Cleanup(cleanup)
	var want int16 = 2357

	if err := buffer.WriteInt16(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadInt16()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Int16 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Uint16(t *testing.T) {
	t.Cleanup(cleanup)
	var want uint16 = 2357

	if err := buffer.WriteUint16(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadUint16()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Uint16 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Int32(t *testing.T) {
	t.Cleanup(cleanup)
	var want int32 = 23571113

	if err := buffer.WriteInt32(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadInt32()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Int32 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Uint32(t *testing.T) {
	t.Cleanup(cleanup)
	var want uint32 = 23571113

	if err := buffer.WriteUint32(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadUint32()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Uint32 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Int64(t *testing.T) {
	t.Cleanup(cleanup)
	var want int64 = 235711131719232931

	if err := buffer.WriteInt64(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadInt64()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Int64 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Uint64(t *testing.T) {
	t.Cleanup(cleanup)
	var want uint64 = 235711131719232931

	if err := buffer.WriteUint64(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadUint64()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Uint64 was incorrect, got: %d, want: %d.", got, want)
	}
}

func TestBuffer_Float32(t *testing.T) {
	t.Cleanup(cleanup)
	var want float32 = 3.1415926535

	if err := buffer.WriteFloat32(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadFloat32()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Float32 was incorrect, got: %f, want: %f.", got, want)
	}
}

func TestBuffer_Float64(t *testing.T) {
	t.Cleanup(cleanup)
	var want = 1.4142135623

	if err := buffer.WriteFloat64(want); err != nil {
		t.Fatal(err)
	}

	got, err := buffer.ReadFloat64()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("Float64 was incorrect, got: %f, want: %f.", got, want)
	}
}
