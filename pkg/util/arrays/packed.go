package arrays

import "math"

type (
	PackedArray interface {
		GetBacking() []int
		GetCapacity() int
		GetBitsPerValue() int
		GetLargestPossibleValue() int
		Set(index, value int)
		Get(index int) int
	}

	packedArray struct {
		backing      []int
		capacity     int
		bitsPerValue int
		valueMask    int
	}
)

func (array *packedArray) GetBacking() []int {
	return array.backing
}

func (array *packedArray) GetCapacity() int {
	return array.capacity
}

func (array *packedArray) GetBitsPerValue() int {
	return array.bitsPerValue
}

func (array *packedArray) GetLargestPossibleValue() int {
	return array.valueMask
}

func (array *packedArray) Set(index, value int) {
	index *= array.bitsPerValue
	i0 := index >> 6
	i1 := index & 0x3f

	array.backing[i0] = array.backing[i0] & ^(array.valueMask<<i1) | (value&array.valueMask)<<i1
	i2 := i1 + array.bitsPerValue
	if i2 > 64 {
		i0++
		array.backing[i0] = array.backing[i0] & ^((1<<i2-64)-1) | value>>(64-i1)
	}
}

func (array *packedArray) Get(index int) int {
	index *= array.bitsPerValue
	i0 := index >> 6
	i1 := index & 0x3f

	value := array.backing[i0] >> i1
	i2 := i1 + array.bitsPerValue
	if i2 > 64 {
		i0++
		value |= array.backing[i0] << (64 - i1)
	}

	return value & array.valueMask
}

func NewPackedArray(bitsPerValue, capacity int) PackedArray {
	return &packedArray{
		backing:      make([]int, int(math.Ceil(float64((bitsPerValue*capacity)/64)))),
		capacity:     capacity,
		bitsPerValue: bitsPerValue,
		valueMask:    (1 << bitsPerValue) - 1,
	}
}
