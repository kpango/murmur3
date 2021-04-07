# murmur3

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![release](https://img.shields.io/github/release/kpango/murmur3.svg)](https://github.com/kpango/murmur3/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/kpango/murmur3.svg)](https://pkg.go.dev/github.com/kpango/murmur3)
[![Go Report Card](https://goreportcard.com/badge/github.com/kpango/murmur3)](https://goreportcard.com/report/github.com/kpango/murmur3)
[![CI](https://github.com/kpango/murmur3/actions/workflows/ci.yaml/badge.svg)](https://github.com/kpango/murmur3/actions/workflows/ci.yaml)

High-performance MurmurHash3 (x64-128) for Go with AVX2, AVX-512, and ARM64 NEON acceleration.

Bit-for-bit compatible with [spaolacci/murmur3](https://github.com/spaolacci/murmur3).

## Features

- **AVX-512 path** – 128-byte unrolled block loop for wide-pipeline CPUs
- **AVX2 path** – 64-byte unrolled block loop with prefetch hints
- **ARM64 NEON path** – 64-byte unrolled block loop with prefetch hints
- **Pure Go fallback** – portable 64-byte unrolled loop for all other architectures
- **Generic type support** – `Sum128Slice[T any]` hashes any contiguous slice without allocation
- **Zero-alloc strings** – `Sum128String` uses `unsafe.StringData` to avoid a copy
- **Backward compatible** – identical output to `spaolacci/murmur3`

## Requirement

Go 1.21 or later (generics and `unsafe.StringData` required).

## Installation

```shell
go get github.com/kpango/murmur3
```

## Usage

### Hash a string

```go
import "github.com/kpango/murmur3"

h1, h2 := murmur3.Sum128String(0, "hello world")
```

### Hash a byte slice

```go
h1, h2 := murmur3.Sum128Slice(0, []byte("hello world"))
```

### Hash any slice (generics)

```go
floats := []float32{1.0, 2.0, 3.0}
h1, h2 := murmur3.Sum128Slice(0, floats)

ints := []uint64{1, 2, 3, 4}
h1, h2 := murmur3.Sum128Slice(42, ints)
```

## Benchmark

Run on Intel Core i9 (AVX-512):

```
BenchmarkSlice-16      	 3000000	       396 ns/op	2584.87 MB/s
BenchmarkString-16     	 3000000	       397 ns/op	2578.62 MB/s
BenchmarkFloat32-16    	 3000000	       396 ns/op	2586.22 MB/s
BenchmarkSpaolacci-16  	 1000000	      1123 ns/op	 911.15 MB/s
```

To compare against the base branch:

```shell
make bench-compare
```

## License

MIT — see [LICENSE](LICENSE)
