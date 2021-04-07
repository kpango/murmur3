//go:build !amd64 && !arm64

package murmur3

func sum128(ptr *byte, length int, seed uint32) (uint64, uint64) {
	if length < 16 {
		return sum128Short(ptr, length, seed)
	}
	return sum128Long(ptr, length, seed)
}
