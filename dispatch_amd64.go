package murmur3

import "golang.org/x/sys/cpu"

var (
	hasAVX512 = cpu.X86.HasAVX512F
	hasAVX2   = cpu.X86.HasAVX2
)

//go:inline
func sum128(ptr *byte, length int, seed uint32) (uint64, uint64) {
	if length < 16 {
		return sum128Short(ptr, length, seed)
	}
	if hasAVX512 {
		return sum128AVX512(ptr, length, seed)
	}
	if hasAVX2 {
		return sum128AVX2(ptr, length, seed)
	}
	return sum128Long(ptr, length, seed)
}
