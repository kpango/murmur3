package murmur3

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2
	hasNEON = cpu.ARM64.HasASIMD
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

// sum128 is the internal dispatcher.
func sum128(ptr *byte, length int, seed uint32) (uint64, uint64) {
	if length < 16 {
		return sum128Short(ptr, length, seed)
	}
	if hasAVX2 {
		h1, h2 := sum128AVX2(ptr, length, seed)
		rem := length & 15
		if rem > 0 {
			k1, k2 := tail(ptr, rem, length-rem)
			h1 ^= k1
			h2 ^= k2
		}
		return final(h1, h2, length)
	}
	return sum128Long(ptr, length, seed)
}
