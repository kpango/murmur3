package main

import (
	"fmt"

	"github.com/kpango/murmur3"
)

func main() {
	fmt.Println(fmix64(10))
	fmt.Println(rotl64(10, 8))
	h1, h2 := murmur3.New(123456).Sum128([]byte("helloqwertyuiophelloqwertyuiop"))

	fmt.Println(h1)
	fmt.Println(h2)
}

func fmix64(x uint64) uint64 {
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33
	return x
}

func rotl64(x uint64, r uint8) uint64 {
	return (x << r) | (x >> (64 - r))
}
