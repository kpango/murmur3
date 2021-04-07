package murmur3

// sum128NEON is an ARM64 NEON-accelerated version of Murmur3 128-bit block processing.
func sum128NEON(ptr *byte, length int, seed uint32) (uint64, uint64)
