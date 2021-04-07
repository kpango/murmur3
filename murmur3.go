package murmur3

import (
	"unsafe"
)

// Sum128String computes the 128-bit MurmurHash3 of a string.
// It avoids allocation by using unsafe.StringData.
func Sum128String(seed uint32, s string) (uint64, uint64) {
	if len(s) == 0 {
		return final(uint64(seed), uint64(seed), 0)
	}
	ptr := unsafe.StringData(s)
	return sum128(ptr, len(s), seed)
}

// Sum128Slice computes the 128-bit MurmurHash3 of a generic slice.
// It avoids allocation by casting the slice directly to a byte pointer.
func Sum128Slice[T any](seed uint32, slice []T) (uint64, uint64) {
	if len(slice) == 0 {
		return final(uint64(seed), uint64(seed), 0)
	}
	var t T
	sizeOfT := int(unsafe.Sizeof(t))
	byteLen := len(slice) * sizeOfT
	ptr := unsafe.Pointer(&slice[0])
	return sum128((*byte)(ptr), byteLen, seed)
}
