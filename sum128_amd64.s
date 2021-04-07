#include "textflag.h"

// Constant values for murmur3 128 bit
#define C1 $0x87c37b91114253d5
#define C2 $0x4cf5ad432745937f

// ONE_BLOCK processes 16 bytes at R10+off and R10+off+8.
// k1/k2 use DX/BX; h1/h2 use AX/CX; c1/c2 use R8/R9.
#define ONE_BLOCK(off) \
	MOVQ  off(R10), DX               \
	MOVQ  (off+8)(R10), BX           \
	IMULQ R8, DX                     \
	IMULQ R9, BX                     \
	ROLQ  $31, DX                    \
	ROLQ  $33, BX                    \
	IMULQ R9, DX                     \
	IMULQ R8, BX                     \
	XORQ  DX, AX                     \
	ROLQ  $27, AX                    \
	ADDQ  CX, AX                     \
	XORQ  BX, CX                     \
	ROLQ  $31, CX                    \
	LEAQ  0x52dce729(AX)(AX*4), AX   \
	ADDQ  AX, CX                     \
	LEAQ  0x38495ab5(CX)(CX*4), CX

// sum128core<>: 4x-unrolled MurmurHash3 x64-128 core, register-calling convention.
//
// Entry register state:
//   R10 = data ptr
//   R11 = data len  (clobbered during tail dispatch)
//   AX  = h1 (= uint64(seed))
//   CX  = h2 (= uint64(seed))
//   R8  = c1
//   R9  = c2
//   R13 = &ret[0]  (16-byte return slot: h1 at 0, h2 at 8)
//
// R12 is saved here as original length for the finalize XOR.
// DX, BX are clobbered (k1/k2 scratch and fmix64 scratch).
TEXT sum128core<>(SB), NOSPLIT, $0
	MOVQ R11, R12           // R12 = original length

	// 4x-unrolled main loop (64 bytes / iter)
core_loop64:
	CMPQ R11, $0x40
	JL   core_loop16

	ONE_BLOCK(0)
	ONE_BLOCK(16)
	ONE_BLOCK(32)
	ONE_BLOCK(48)

	ADDQ $0x40, R10
	SUBQ $0x40, R11
	JMP  core_loop64

core_loop16:
	CMPQ R11, $0x10
	JL   core_do_tail

	ONE_BLOCK(0)

	ADDQ $0x10, R10
	SUBQ $0x10, R11
	JMP  core_loop16

	// Tail: 0-15 remaining bytes.
	// Control flow:
	//   0 bytes    → core_finalize
	//   1-7 bytes  → core_k1_only → core_k1_done → core_finalize
	//   8 bytes    → core_exactly_8 → core_finalize
	//   9-15 bytes → k2 dispatch → core_k2_done → core_exactly_8 → core_finalize
core_do_tail:
	TESTQ R11, R11
	JZ    core_finalize

	CMPQ  R11, $8
	JL    core_k1_only
	JE    core_exactly_8

	// 9-15 bytes: process k2 (bytes [8..]) then fall through to k1 (bytes [0..7])
	SUBQ  $8, R11

	CMPQ  R11, $4
	JGE   core_k2_4_7
	CMPQ  R11, $2
	JGE   core_k2_2_3
	MOVBQZX 8(R10), BX
	JMP   core_k2_done

core_k2_2_3:
	JE    core_k2_2
	MOVWQZX 8(R10), BX
	MOVBQZX 10(R10), DX
	SHLQ    $16, DX
	ORQ     DX, BX
	JMP   core_k2_done
core_k2_2:
	MOVWQZX 8(R10), BX
	JMP   core_k2_done

core_k2_4_7:
	JE    core_k2_4
	CMPQ  R11, $6
	JGE   core_k2_6_7
	MOVLQZX 8(R10), BX
	MOVBQZX 12(R10), DX
	SHLQ    $32, DX
	ORQ     DX, BX
	JMP   core_k2_done
core_k2_6_7:
	JE    core_k2_6
	MOVLQZX 8(R10), BX
	MOVWQZX 12(R10), DX
	SHLQ    $32, DX
	ORQ     DX, BX
	MOVBQZX 14(R10), DX
	SHLQ    $48, DX
	ORQ     DX, BX
	JMP   core_k2_done
core_k2_6:
	MOVLQZX 8(R10), BX
	MOVWQZX 12(R10), DX
	SHLQ    $32, DX
	ORQ     DX, BX
	JMP   core_k2_done
core_k2_4:
	MOVLQZX 8(R10), BX

core_k2_done:
	IMULQ R9, BX
	ROLQ  $33, BX
	IMULQ R8, BX
	XORQ  BX, CX

core_exactly_8:
	MOVQ  (R10), DX
	IMULQ R8, DX
	ROLQ  $31, DX
	IMULQ R9, DX
	XORQ  DX, AX
	JMP   core_finalize

core_k1_only:
	CMPQ  R11, $4
	JGE   core_k1_4_7
	CMPQ  R11, $2
	JGE   core_k1_2_3
	MOVBQZX (R10), DX
	JMP   core_k1_done

core_k1_2_3:
	JE    core_k1_2
	MOVWQZX (R10), DX
	MOVBQZX 2(R10), BX
	SHLQ    $16, BX
	ORQ     BX, DX
	JMP   core_k1_done
core_k1_2:
	MOVWQZX (R10), DX
	JMP   core_k1_done

core_k1_4_7:
	JE    core_k1_4
	CMPQ  R11, $6
	JGE   core_k1_6_7
	MOVLQZX (R10), DX
	MOVBQZX 4(R10), BX
	SHLQ    $32, BX
	ORQ     BX, DX
	JMP   core_k1_done
core_k1_6_7:
	JE    core_k1_6
	MOVLQZX (R10), DX
	MOVWQZX 4(R10), BX
	SHLQ    $32, BX
	ORQ     BX, DX
	MOVBQZX 6(R10), BX
	SHLQ    $48, BX
	ORQ     BX, DX
	JMP   core_k1_done
core_k1_6:
	MOVLQZX (R10), DX
	MOVWQZX 4(R10), BX
	SHLQ    $32, BX
	ORQ     BX, DX
	JMP   core_k1_done
core_k1_4:
	MOVLQZX (R10), DX

core_k1_done:
	IMULQ R8, DX
	ROLQ  $31, DX
	IMULQ R9, DX
	XORQ  DX, AX

core_finalize:
	XORQ  R12, AX           // h1 ^= original_len
	XORQ  R12, CX           // h2 ^= original_len

	ADDQ  CX, AX            // h1 += h2
	ADDQ  AX, CX            // h2 += h1

	// Interleaved fmix64: process h1 (AX) and h2 (CX) in parallel.
	// Round 1
	MOVQ  AX, DX
	MOVQ  CX, BX
	SHRQ  $33, DX
	SHRQ  $33, BX
	XORQ  DX, AX
	XORQ  BX, CX

	MOVQ  $0xff51afd7ed558ccd, DX
	IMULQ DX, AX
	IMULQ DX, CX

	// Round 2
	MOVQ  AX, DX
	MOVQ  CX, BX
	SHRQ  $33, DX
	SHRQ  $33, BX
	XORQ  DX, AX
	XORQ  BX, CX

	MOVQ  $0xc4ceb9fe1a85ec53, DX
	IMULQ DX, AX
	IMULQ DX, CX

	// Round 3
	MOVQ  AX, DX
	MOVQ  CX, BX
	SHRQ  $33, DX
	SHRQ  $33, BX
	XORQ  DX, AX
	XORQ  BX, CX

	ADDQ  CX, AX            // h1 += h2
	ADDQ  AX, CX            // h2 += h1

	MOVQ  AX, (R13)
	MOVQ  CX, 8(R13)
	RET

// sum128core512<>: 8x-unrolled MurmurHash3 x64-128 core.
// Same register convention as sum128core<>.
TEXT sum128core512<>(SB), NOSPLIT, $0
	MOVQ R11, R12

	// 8x-unrolled main loop (128 bytes / iter)
c512_loop128:
	CMPQ R11, $0x80
	JL   c512_loop64

	ONE_BLOCK(0)
	ONE_BLOCK(16)
	ONE_BLOCK(32)
	ONE_BLOCK(48)
	ONE_BLOCK(64)
	ONE_BLOCK(80)
	ONE_BLOCK(96)
	ONE_BLOCK(112)

	ADDQ $0x80, R10
	SUBQ $0x80, R11
	JMP  c512_loop128

	// 4x-unrolled fallback (64 bytes / iter)
c512_loop64:
	CMPQ R11, $0x40
	JL   c512_loop16

	ONE_BLOCK(0)
	ONE_BLOCK(16)
	ONE_BLOCK(32)
	ONE_BLOCK(48)

	ADDQ $0x40, R10
	SUBQ $0x40, R11
	JMP  c512_loop64

c512_loop16:
	CMPQ R11, $0x10
	JL   c512_do_tail

	ONE_BLOCK(0)

	ADDQ $0x10, R10
	SUBQ $0x10, R11
	JMP  c512_loop16

c512_do_tail:
	TESTQ R11, R11
	JZ    c512_finalize

	CMPQ  R11, $8
	JL    c512_k1_only
	JE    c512_exactly_8

	SUBQ  $8, R11

	CMPQ  R11, $4
	JGE   c512_k2_4_7
	CMPQ  R11, $2
	JGE   c512_k2_2_3
	MOVBQZX 8(R10), BX
	JMP   c512_k2_done

c512_k2_2_3:
	JE    c512_k2_2
	MOVWQZX 8(R10), BX
	MOVBQZX 10(R10), DX
	SHLQ    $16, DX
	ORQ     DX, BX
	JMP   c512_k2_done
c512_k2_2:
	MOVWQZX 8(R10), BX
	JMP   c512_k2_done

c512_k2_4_7:
	JE    c512_k2_4
	CMPQ  R11, $6
	JGE   c512_k2_6_7
	MOVLQZX 8(R10), BX
	MOVBQZX 12(R10), DX
	SHLQ    $32, DX
	ORQ     DX, BX
	JMP   c512_k2_done
c512_k2_6_7:
	JE    c512_k2_6
	MOVLQZX 8(R10), BX
	MOVWQZX 12(R10), DX
	SHLQ    $32, DX
	ORQ     DX, BX
	MOVBQZX 14(R10), DX
	SHLQ    $48, DX
	ORQ     DX, BX
	JMP   c512_k2_done
c512_k2_6:
	MOVLQZX 8(R10), BX
	MOVWQZX 12(R10), DX
	SHLQ    $32, DX
	ORQ     DX, BX
	JMP   c512_k2_done
c512_k2_4:
	MOVLQZX 8(R10), BX

c512_k2_done:
	IMULQ R9, BX
	ROLQ  $33, BX
	IMULQ R8, BX
	XORQ  BX, CX

c512_exactly_8:
	MOVQ  (R10), DX
	IMULQ R8, DX
	ROLQ  $31, DX
	IMULQ R9, DX
	XORQ  DX, AX
	JMP   c512_finalize

c512_k1_only:
	CMPQ  R11, $4
	JGE   c512_k1_4_7
	CMPQ  R11, $2
	JGE   c512_k1_2_3
	MOVBQZX (R10), DX
	JMP   c512_k1_done

c512_k1_2_3:
	JE    c512_k1_2
	MOVWQZX (R10), DX
	MOVBQZX 2(R10), BX
	SHLQ    $16, BX
	ORQ     BX, DX
	JMP   c512_k1_done
c512_k1_2:
	MOVWQZX (R10), DX
	JMP   c512_k1_done

c512_k1_4_7:
	JE    c512_k1_4
	CMPQ  R11, $6
	JGE   c512_k1_6_7
	MOVLQZX (R10), DX
	MOVBQZX 4(R10), BX
	SHLQ    $32, BX
	ORQ     BX, DX
	JMP   c512_k1_done
c512_k1_6_7:
	JE    c512_k1_6
	MOVLQZX (R10), DX
	MOVWQZX 4(R10), BX
	SHLQ    $32, BX
	ORQ     BX, DX
	MOVBQZX 6(R10), BX
	SHLQ    $48, BX
	ORQ     BX, DX
	JMP   c512_k1_done
c512_k1_6:
	MOVLQZX (R10), DX
	MOVWQZX 4(R10), BX
	SHLQ    $32, BX
	ORQ     BX, DX
	JMP   c512_k1_done
c512_k1_4:
	MOVLQZX (R10), DX

c512_k1_done:
	IMULQ R8, DX
	ROLQ  $31, DX
	IMULQ R9, DX
	XORQ  DX, AX

c512_finalize:
	XORQ  R12, AX
	XORQ  R12, CX

	ADDQ  CX, AX
	ADDQ  AX, CX

	MOVQ  AX, DX
	MOVQ  CX, BX
	SHRQ  $33, DX
	SHRQ  $33, BX
	XORQ  DX, AX
	XORQ  BX, CX

	MOVQ  $0xff51afd7ed558ccd, DX
	IMULQ DX, AX
	IMULQ DX, CX

	MOVQ  AX, DX
	MOVQ  CX, BX
	SHRQ  $33, DX
	SHRQ  $33, BX
	XORQ  DX, AX
	XORQ  BX, CX

	MOVQ  $0xc4ceb9fe1a85ec53, DX
	IMULQ DX, AX
	IMULQ DX, CX

	MOVQ  AX, DX
	MOVQ  CX, BX
	SHRQ  $33, DX
	SHRQ  $33, BX
	XORQ  DX, AX
	XORQ  BX, CX

	ADDQ  CX, AX
	ADDQ  AX, CX

	MOVQ  AX, (R13)
	MOVQ  CX, 8(R13)
	RET

// func sum128AVX2(ptr *byte, length int, seed uint32) (uint64, uint64)
//
// ABI0 frame: ptr(8)+length(8)+seed(4)+_pad(4)+h1(8)+h2(8) = 40 bytes
// Called only when length >= 16 (dispatch_amd64.go guarantees this).
TEXT ·sum128AVX2(SB), NOSPLIT, $0-40
	MOVQ ptr+0(FP), R10
	MOVQ length+8(FP), R11
	MOVL seed+16(FP), AX
	MOVQ AX, CX
	MOVQ C1, R8
	MOVQ C2, R9
	LEAQ ret+24(FP), R13
	JMP  sum128core<>(SB)

// func sum128AVX512(ptr *byte, length int, seed uint32) (uint64, uint64)
//
// ABI0 frame: ptr(8)+length(8)+seed(4)+_pad(4)+h1(8)+h2(8) = 40 bytes
// Called only when length >= 16 (dispatch_amd64.go guarantees this).
TEXT ·sum128AVX512(SB), NOSPLIT, $0-40
	MOVQ ptr+0(FP), R10
	MOVQ length+8(FP), R11
	MOVL seed+16(FP), AX
	MOVQ AX, CX
	MOVQ C1, R8
	MOVQ C2, R9
	LEAQ ret+24(FP), R13
	JMP  sum128core512<>(SB)

// func Sum128(seed uint32, data []byte) (uint64, uint64)
//
// ABI0 frame layout (seed uint32 aligns data.ptr to +8):
//   +0:  seed   (uint32, 4 bytes)
//   +4:  _pad   (4 bytes)
//   +8:  data.ptr (8 bytes)
//   +16: data.len (8 bytes)
//   +24: data.cap (8 bytes)
//   +32: ret h1  (8 bytes)
//   +40: ret h2  (8 bytes)
//        total: 48 bytes
TEXT ·Sum128(SB), NOSPLIT, $0-48
	MOVQ data_base+8(FP), R10
	MOVQ data_len+16(FP), R11
	MOVL seed+0(FP), AX
	MOVQ AX, CX
	MOVQ C1, R8
	MOVQ C2, R9
	LEAQ ret+32(FP), R13
	JMP  sum128core<>(SB)
