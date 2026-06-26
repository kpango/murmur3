package main

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

func main() {
	build.TEXT("sum128AVX2", build.NOSPLIT, "func(ptr *byte, length int, seed uint32) (uint64, uint64)")
	build.Doc("sum128AVX2 is an AVX2 accelerated version of Murmur3 128-bit block processing.")

	ptr := build.Load(build.Param("ptr"), reg.R10)
	length := build.Load(build.Param("length"), reg.R11)

	seed64 := reg.R12
	build.XORQ(seed64, seed64)
	build.Load(build.Param("seed"), seed64.As32())

	c1 := reg.R8
	build.MOVQ(operand.U64(0x87c37b91114253d5), c1)

	c2 := reg.R9
	build.MOVQ(operand.U64(0x4cf5ad432745937f), c2)

	h1 := reg.RAX
	h2 := reg.RCX
	build.MOVQ(seed64, h1)
	build.MOVQ(seed64, h2)

	build.CMPQ(length, operand.U32(64))
	build.JL(operand.LabelRef("loop16"))

	build.Label("loop64")
	for i := 0; i < 4; i++ {
		k1 := reg.RDX
		k2 := reg.RBX
		build.MOVQ(operand.Mem{Base: ptr, Disp: i * 16}, k1)
		build.MOVQ(operand.Mem{Base: ptr, Disp: i*16 + 8}, k2)

		build.IMULQ(c1, k1)
		build.ROLQ(operand.U8(31), k1)
		build.IMULQ(c2, k1)
		build.XORQ(k1, h1)

		build.ROLQ(operand.U8(27), h1)
		build.ADDQ(h2, h1)
		build.IMUL3Q(operand.U32(5), h1, h1)
		build.ADDQ(operand.U32(0x52dce729), h1)

		build.IMULQ(c2, k2)
		build.ROLQ(operand.U8(33), k2)
		build.IMULQ(c1, k2)
		build.XORQ(k2, h2)

		build.ROLQ(operand.U8(31), h2)
		build.ADDQ(h1, h2)
		build.IMUL3Q(operand.U32(5), h2, h2)
		build.ADDQ(operand.U32(0x38495ab5), h2)
	}

	build.ADDQ(operand.U32(64), ptr)
	build.SUBQ(operand.U32(64), length)
	build.CMPQ(length, operand.U32(64))
	build.JGE(operand.LabelRef("loop64"))

	build.Label("loop16")
	build.CMPQ(length, operand.U32(16))
	build.JL(operand.LabelRef("tail_only"))

	k1 := reg.RDX
	k2 := reg.RBX
	build.MOVQ(operand.Mem{Base: ptr, Disp: 0}, k1)
	build.MOVQ(operand.Mem{Base: ptr, Disp: 8}, k2)

	build.IMULQ(c1, k1)
	build.ROLQ(operand.U8(31), k1)
	build.IMULQ(c2, k1)
	build.XORQ(k1, h1)

	build.ROLQ(operand.U8(27), h1)
	build.ADDQ(h2, h1)
	build.IMUL3Q(operand.U32(5), h1, h1)
	build.ADDQ(operand.U32(0x52dce729), h1)

	build.IMULQ(c2, k2)
	build.ROLQ(operand.U8(33), k2)
	build.IMULQ(c1, k2)
	build.XORQ(k2, h2)

	build.ROLQ(operand.U8(31), h2)
	build.ADDQ(h1, h2)
	build.IMUL3Q(operand.U32(5), h2, h2)
	build.ADDQ(operand.U32(0x38495ab5), h2)

	build.ADDQ(operand.U32(16), ptr)
	build.SUBQ(operand.U32(16), length)
	build.JMP(operand.LabelRef("loop16"))

	build.Label("tail_only")
	build.Store(h1, build.ReturnIndex(0))
	build.Store(h2, build.ReturnIndex(1))

	build.RET()
	build.Generate()
}
