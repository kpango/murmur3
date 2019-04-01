package murmur3

import (
	"unsafe"
)

type MurmurHash3 struct {
	k1 uint64
	k2 uint64
	h1 uint64
	h2 uint64
}

const (
	c1 uint64 = 0x87c37b91114253d5
	c2 uint64 = 0x4cf5ad432745937f
)

func New(seed uint64) *MurmurHash3 {
	return &MurmurHash3{
		h1: seed,
		h2: seed,
	}
}

func (m *MurmurHash3) Sum128(data []byte) (uint64, uint64) {

	var (
		blocks  = *(*[]uint64)(unsafe.Pointer(&data))
		ld      = len(data)
		nblocks = ld >> 4
	)

	m.k1 = 0
	m.k2 = 0
	for i := 0; i < nblocks; i++ {
		m.k1 = blocks[i<<1+0]
		m.k2 = blocks[i<<1+1]

		m.k1Calc()
		m.h1 = (rotl64(m.h1, 27)+m.h2)*5 + 0x52dce729

		m.k2Calc()
		m.h2 = (rotl64(m.h2, 31)+m.h1)*5 + 0x38495ab5
	}

	m.k1 = 0
	m.k2 = 0

	tail := data[nblocks<<4:]
	switch len(tail) & 15 {
	case 15:
		m.k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		m.k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		m.k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		m.k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		m.k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		m.k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		m.k2 ^= uint64(tail[8]) << 0
		m.k2Calc()
		fallthrough
	case 8:
		m.k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		m.k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		m.k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		m.k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		m.k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		m.k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		m.k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		m.k1 ^= uint64(tail[0]) << 0
		m.k1Calc()
	}
	m.h1 ^= *(*uint64)(unsafe.Pointer(&ld))
	m.h2 ^= *(*uint64)(unsafe.Pointer(&ld))

	f1 := fmix64(m.h1 + m.h2)
	f2 := fmix64(m.h1 + m.h2<<1)

	return f1 + f2, f1 + f2<<1
}

func (m *MurmurHash3) k1Calc() {
	m.k1 = rotl64(m.k1*c1, 31) * c2
	m.h1 ^= m.k1
}

func (m *MurmurHash3) k2Calc() {
	m.k2 = rotl64(m.k2*c2, 33) * c1
	m.h2 ^= m.k2
}

func rotl64(x uint64, r uint8) uint64 {
	return (x << r) | (x >> (64 - r))
}

func fmix64(x uint64) uint64 {
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33
	return x
}
