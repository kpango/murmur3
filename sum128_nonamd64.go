//go:build !amd64

package murmur3

// Sum128 computes the 128-bit MurmurHash3 of a byte slice.
func Sum128(seed uint32, data []byte) (uint64, uint64) {
	if len(data) == 0 {
		return final(uint64(seed), uint64(seed), 0)
	}
	return sum128(&data[0], len(data), seed)
}
