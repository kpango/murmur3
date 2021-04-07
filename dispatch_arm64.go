package murmur3

import "golang.org/x/sys/cpu"

var hasNEON = cpu.ARM64.HasASIMD

func sum128(ptr *byte, length int, seed uint32) (uint64, uint64) {
	if length < 16 {
		return sum128Short(ptr, length, seed)
	}
	if hasNEON {
		h1, h2 := sum128NEON(ptr, length, seed)
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
