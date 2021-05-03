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
	startLong := (index * array.bitsPerValue) / 64
	startOffset := (index * array.bitsPerValue) % 64
	endLong := ((index+1)*array.bitsPerValue - 1) / 64

	uValue := uint64(value) & array.valueMask

	array.data[startLong] |= uValue << startOffset
	if startLong != endLong {
		array.data[endLong] = uValue >> (64 - startOffset)
	}
}

func (array *packedArray) Get(index int) int {
	startLong := (index * array.bitsPerValue) / 64
	startOffset := (index * array.bitsPerValue) % 64
	endLong := ((index+1)*array.bitsPerValue - 1) / 64

	var value uint64
	if startLong == endLong {
		value = array.data[startLong] >> startOffset
	} else {
		endOffset := 64 - startOffset
		value = array.data[startLong]>>startOffset | array.data[endLong]<<endOffset
	}
	return int(value & array.valueMask)
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
		data:         make([]uint64, int(math.Ceil(float64((capacity*bitsPerValue)/64)))),
		capacity:     capacity,
		bitsPerValue: bitsPerValue,
		valueMask:    uint64((1 << bitsPerValue) - 1),
	}
}
