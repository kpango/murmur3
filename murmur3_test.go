package murmur3

import (
	"testing"

	ref "github.com/spaolacci/murmur3"
)

func TestCorrectness(t *testing.T) {
	tests := []int{0, 1, 2, 3, 4, 7, 8, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129, 256, 1024, 4096}
	for _, l := range tests {
		b := make([]byte, l)
		for i := range b {
			b[i] = byte(i % 256)
		}

		seed := uint32(0)
		h1Ref, h2Ref := ref.Sum128WithSeed(b, seed)

		h1Act, h2Act := Sum128Slice(seed, b)
		if h1Ref != h1Act || h2Ref != h2Act {
			t.Errorf("Mismatch slice len %d: ref (%x, %x), act (%x, %x)", l, h1Ref, h2Ref, h1Act, h2Act)
		}

		h1ActS, h2ActS := Sum128String(seed, string(b))
		if h1Ref != h1ActS || h2Ref != h2ActS {
			t.Errorf("Mismatch string len %d: ref (%x, %x), act (%x, %x)", l, h1Ref, h2Ref, h1ActS, h2ActS)
		}
	}
}

func FuzzMurmur3(f *testing.F) {
	f.Add([]byte(""), uint32(0))
	f.Add([]byte("hello"), uint32(42))
	f.Add([]byte("1234567890123456"), uint32(0))
	f.Add([]byte("12345678901234567"), uint32(1))

	longInp := make([]byte, 1000)
	for i := range longInp {
		longInp[i] = byte(i)
	}
	f.Add(longInp, uint32(123))

	f.Fuzz(func(t *testing.T, b []byte, seed uint32) {
		h1Ref, h2Ref := ref.Sum128WithSeed(b, seed)
		h1Act, h2Act := Sum128Slice(seed, b)

		if h1Ref != h1Act || h2Ref != h2Act {
			t.Errorf("Slice Mismatch: seed=%d, len=%d, ref=(%x, %x), act=(%x, %x)", seed, len(b), h1Ref, h2Ref, h1Act, h2Act)
		}

		h1ActS, h2ActS := Sum128String(seed, string(b))
		if h1Ref != h1ActS || h2Ref != h2ActS {
			t.Errorf("String Mismatch: seed=%d, len=%d, ref=(%x, %x), act=(%x, %x)", seed, len(b), h1Ref, h2Ref, h1ActS, h2ActS)
		}
	})
}

