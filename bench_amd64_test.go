package murmur3

import (
	"fmt"
	"testing"
)

// BenchmarkAVX512 benchmarks the complete sum128AVX512 function (includes tail + finalization).
func BenchmarkAVX512(b *testing.B) {
	for _, size := range benchSizes {
		data := make([]byte, size)
		ptr := &data[0]
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sum128AVX512(ptr, size, 0)
			}
		})
	}
}

// BenchmarkAVX2 benchmarks the complete sum128AVX2 function (includes tail + finalization).
func BenchmarkAVX2(b *testing.B) {
	for _, size := range benchSizes {
		data := make([]byte, size)
		ptr := &data[0]
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sum128AVX2(ptr, size, 0)
			}
		})
	}
}

// BenchmarkPureGo benchmarks the pure-Go sum128Long function (includes tail + finalization).
func BenchmarkPureGo(b *testing.B) {
	for _, size := range benchSizes {
		data := make([]byte, size)
		ptr := &data[0]
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sum128Long(ptr, size, 0)
			}
		})
	}
}
