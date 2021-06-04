package bytes

import "math"

type (
	PackedArray interface {
		GetData() []uint64
		GetCapacity() int
		GetBitsPerValue() int
		GetValueMask() uint64
		Set(index, value int)
		Get(index int) int
		Resized(bitsPerValue int) PackedArray
	}

	packedArray struct {
		data         []uint64
		capacity     int
		bitsPerValue int
		valueMask    uint64
	}
)

func (array *packedArray) GetData() []uint64 {
	return array.data
}

func (array *packedArray) GetCapacity() int {
	return array.capacity
}

func (array *packedArray) GetBitsPerValue() int {
	return array.bitsPerValue
}

func (array *packedArray) GetValueMask() uint64 {
	return array.valueMask
}

func (array *packedArray) Set(index, value int) {
	long := index / (64 / array.bitsPerValue)
	offset := (index - long*(64/array.bitsPerValue)) * array.bitsPerValue
	array.data[long] = array.data[long]&(array.valueMask<<offset^math.MaxUint64) | (uint64(value)&array.valueMask)<<offset
}

func (array *packedArray) Get(index int) int {
	long := index / (64 / array.bitsPerValue)
	offset := (index - long*(64/array.bitsPerValue)) * array.bitsPerValue
	return int(array.data[long] >> offset & array.valueMask)
}

func (array *packedArray) Resized(bitsPerValue int) PackedArray {
	newArray := NewPackedArray(bitsPerValue, array.capacity)
	for i := 0; i < array.capacity; i++ {
		newArray.Set(i, array.Get(i))
	}
	return newArray
}

func NewPackedArray(bitsPerValue, capacity int) PackedArray {
	return &packedArray{
		data:         make([]uint64, (capacity+(64/bitsPerValue)-1)/(64/bitsPerValue)),
		capacity:     capacity,
		bitsPerValue: bitsPerValue,
		valueMask:    uint64((1 << bitsPerValue) - 1),
	}
}
