package murmur3

import (
	"fmt"
	"testing"

	ref "github.com/spaolacci/murmur3"
	twmb "github.com/twmb/murmur3"
)

var benchSizes = []int{16, 64, 128, 256, 512, 1024, 4096, 65536}

func BenchmarkOursSlice(b *testing.B) {
	for _, size := range benchSizes {
		data := make([]byte, size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Sum128(0, data)
			}
		})
	}
}

func BenchmarkOursString(b *testing.B) {
	for _, size := range benchSizes {
		data := string(make([]byte, size))
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Sum128String(0, data)
			}
		})
	}
}

func BenchmarkSpaolacci(b *testing.B) {
	for _, size := range benchSizes {
		data := make([]byte, size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ref.Sum128WithSeed(data, 0)
			}
		})
	}
}

func BenchmarkTwmb(b *testing.B) {
	for _, size := range benchSizes {
		data := make([]byte, size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				twmb.Sum128(data)
			}
		})
	}
}
