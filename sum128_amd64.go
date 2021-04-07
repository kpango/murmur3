package murmur3

//go:noescape
func Sum128(seed uint32, data []byte) (uint64, uint64)

//go:noescape
func sum128AVX2(ptr *byte, length int, seed uint32) (uint64, uint64)

//go:noescape
func sum128AVX512(ptr *byte, length int, seed uint32) (uint64, uint64)
