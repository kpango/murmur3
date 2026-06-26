package murmur3

import (
	"unsafe"
)

const (
	c1_128 = 0x87c37b91114253d5
	c2_128 = 0x4cf5ad432745937f
)

func bitsRotl64(x uint64, k int) uint64 {
	return (x << k) | (x >> (64 - k))
}

// final applies the finalization steps of MurmurHash3
//
//go:inline
func final(h1, h2 uint64, clen int) (uint64, uint64) {
	h1 ^= uint64(clen)
	h2 ^= uint64(clen)

	h1 += h2
	h2 += h1

	h1 ^= h1 >> 33
	h1 *= 0xff51afd7ed558ccd
	h1 ^= h1 >> 33
	h1 *= 0xc4ceb9fe1a85ec53
	h1 ^= h1 >> 33

	h2 ^= h2 >> 33
	h2 *= 0xff51afd7ed558ccd
	h2 ^= h2 >> 33
	h2 *= 0xc4ceb9fe1a85ec53
	h2 ^= h2 >> 33

	h1 += h2
	h2 += h1

	return h1, h2
}

// sum128Short handles data less than 16 bytes.
// Forced inline.
//
//go:inline
func sum128Short(ptr *byte, length int, seed uint32) (uint64, uint64) {
	h1 := uint64(seed)
	h2 := uint64(seed)

	k1, k2 := tail(ptr, length, 0)
	h1 ^= k1
	h2 ^= k2

	return final(h1, h2, length)
}

// sum128Long handles data of 16 bytes or more.
func sum128Long(ptr *byte, length int, seed uint32) (uint64, uint64) {
	h1 := uint64(seed)
	h2 := uint64(seed)

	offset := 0
	limit := length - 15

	limit64 := length - 63
	for offset < limit64 {
		p := unsafe.Add(unsafe.Pointer(ptr), offset)

		k1 := *(*uint64)(p)
		k2 := *(*uint64)(unsafe.Add(p, 8))
		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> 33)
		k1 *= c2_128
		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> 37)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> 31)
		k2 *= c1_128
		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> 33)
		h2 += h1
		h2 = h2*5 + 0x38495ab5

		k3 := *(*uint64)(unsafe.Add(p, 16))
		k4 := *(*uint64)(unsafe.Add(p, 24))
		k3 *= c1_128
		k3 = (k3 << 31) | (k3 >> 33)
		k3 *= c2_128
		h1 ^= k3
		h1 = (h1 << 27) | (h1 >> 37)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k4 *= c2_128
		k4 = (k4 << 33) | (k4 >> 31)
		k4 *= c1_128
		h2 ^= k4
		h2 = (h2 << 31) | (h2 >> 33)
		h2 += h1
		h2 = h2*5 + 0x38495ab5

		k5 := *(*uint64)(unsafe.Add(p, 32))
		k6 := *(*uint64)(unsafe.Add(p, 40))
		k5 *= c1_128
		k5 = (k5 << 31) | (k5 >> 33)
		k5 *= c2_128
		h1 ^= k5
		h1 = (h1 << 27) | (h1 >> 37)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k6 *= c2_128
		k6 = (k6 << 33) | (k6 >> 31)
		k6 *= c1_128
		h2 ^= k6
		h2 = (h2 << 31) | (h2 >> 33)
		h2 += h1
		h2 = h2*5 + 0x38495ab5

		k7 := *(*uint64)(unsafe.Add(p, 48))
		k8 := *(*uint64)(unsafe.Add(p, 56))
		k7 *= c1_128
		k7 = (k7 << 31) | (k7 >> 33)
		k7 *= c2_128
		h1 ^= k7
		h1 = (h1 << 27) | (h1 >> 37)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k8 *= c2_128
		k8 = (k8 << 33) | (k8 >> 31)
		k8 *= c1_128
		h2 ^= k8
		h2 = (h2 << 31) | (h2 >> 33)
		h2 += h1
		h2 = h2*5 + 0x38495ab5

		offset += 64
	}

	for offset < limit {
		p := unsafe.Add(unsafe.Pointer(ptr), offset)
		k1 := *(*uint64)(p)
		k2 := *(*uint64)(unsafe.Add(p, 8))
		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> 33)
		k1 *= c2_128
		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> 37)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> 31)
		k2 *= c1_128
		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> 33)
		h2 += h1
		h2 = h2*5 + 0x38495ab5

		offset += 16
	}

	tailLen := length - offset
	if tailLen > 0 {
		k1, k2 := tail(ptr, tailLen, offset)
		h1 ^= k1
		h2 ^= k2
	}

	return final(h1, h2, length)
}

// tail processes the remaining 1 to 15 bytes.
// Forced inline.
//
//go:inline
func tail(ptr *byte, length int, offset int) (uint64, uint64) {
	var k1, k2 uint64
	p := unsafe.Add(unsafe.Pointer(ptr), offset)

	switch length {
	case 15:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 14))) << 48
		fallthrough
	case 14:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 13))) << 40
		fallthrough
	case 13:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 12))) << 32
		fallthrough
	case 12:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 11))) << 24
		fallthrough
	case 11:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 10))) << 16
		fallthrough
	case 10:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 9))) << 8
		fallthrough
	case 9:
		k2 ^= uint64(*(*byte)(unsafe.Add(p, 8))) << 0
		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> 31)
		k2 *= c1_128
		fallthrough

	case 8:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 7))) << 56
		fallthrough
	case 7:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 6))) << 48
		fallthrough
	case 6:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 5))) << 40
		fallthrough
	case 5:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 4))) << 32
		fallthrough
	case 4:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 3))) << 24
		fallthrough
	case 3:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 2))) << 16
		fallthrough
	case 2:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 1))) << 8
		fallthrough
	case 1:
		k1 ^= uint64(*(*byte)(unsafe.Add(p, 0))) << 0
		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> 33)
		k1 *= c2_128
	}

	return k1, k2
}
