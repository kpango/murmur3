#include "textflag.h"

// Constant values for murmur3 128 bit
#define C1 $0x87c37b91114253d5
#define C2 $0x4cf5ad432745937f
#define C3 $0x52dce729
#define C4 $0x38495ab5

// func sum128NEON(ptr *byte, length int, seed uint32) (uint64, uint64)
// We will also use scalar optimized unrolling since sequential chaining makes real NEON hard to pipeline fully.
TEXT ·sum128NEON(SB), NOSPLIT, $0-40
	MOVD ptr+0(FP), R0
	MOVD length+8(FP), R1
	MOVWU seed+16(FP), R2

	MOVD C1, R3
	MOVD C2, R4
	MOVD R2, R5 // h1
	MOVD R2, R6 // h2

	CMP $0x40, R1
	BLT loop16

loop64:
	PRFM PLDL1KEEP, 256(R0)
	// Block 1
	MOVD (R0), R7
	MOVD 8(R0), R8

	MUL R3, R7
	ROR $33, R7, R7 // 64-31 = 33
	MUL R4, R7
	EOR R7, R5
	ROR $37, R5, R5 // 64-27 = 37
	ADD R6, R5
	MOVD $5, R9
	MUL R9, R5
	MOVD C3, R10
	ADD R10, R5

	MUL R4, R8
	ROR $31, R8, R8 // 64-33 = 31
	MUL R3, R8
	EOR R8, R6
	ROR $33, R6, R6 // 64-31 = 33
	ADD R5, R6
	MUL R9, R6
	MOVD C4, R10
	ADD R10, R6

	// Block 2
	MOVD 16(R0), R7
	MOVD 24(R0), R8

	MUL R3, R7
	ROR $33, R7, R7
	MUL R4, R7
	EOR R7, R5
	ROR $37, R5, R5
	ADD R6, R5
	MUL R9, R5
	MOVD C3, R10
	ADD R10, R5

	MUL R4, R8
	ROR $31, R8, R8
	MUL R3, R8
	EOR R8, R6
	ROR $33, R6, R6
	ADD R5, R6
	MUL R9, R6
	MOVD C4, R10
	ADD R10, R6

	// Block 3
	MOVD 32(R0), R7
	MOVD 40(R0), R8

	MUL R3, R7
	ROR $33, R7, R7
	MUL R4, R7
	EOR R7, R5
	ROR $37, R5, R5
	ADD R6, R5
	MUL R9, R5
	MOVD C3, R10
	ADD R10, R5

	MUL R4, R8
	ROR $31, R8, R8
	MUL R3, R8
	EOR R8, R6
	ROR $33, R6, R6
	ADD R5, R6
	MUL R9, R6
	MOVD C4, R10
	ADD R10, R6

	// Block 4
	MOVD 48(R0), R7
	MOVD 56(R0), R8

	MUL R3, R7
	ROR $33, R7, R7
	MUL R4, R7
	EOR R7, R5
	ROR $37, R5, R5
	ADD R6, R5
	MUL R9, R5
	MOVD C3, R10
	ADD R10, R5

	MUL R4, R8
	ROR $31, R8, R8
	MUL R3, R8
	EOR R8, R6
	ROR $33, R6, R6
	ADD R5, R6
	MUL R9, R6
	MOVD C4, R10
	ADD R10, R6

	ADD $0x40, R0
	SUB $0x40, R1
	CMP $0x40, R1
	BGE loop64

loop16:
	CMP $0x10, R1
	BLT tail_only

	MOVD (R0), R7
	MOVD 8(R0), R8

	MUL R3, R7
	ROR $33, R7, R7
	MUL R4, R7
	EOR R7, R5
	ROR $37, R5, R5
	ADD R6, R5
	MOVD $5, R9
	MUL R9, R5
	MOVD C3, R10
	ADD R10, R5

	MUL R4, R8
	ROR $31, R8, R8
	MUL R3, R8
	EOR R8, R6
	ROR $33, R6, R6
	ADD R5, R6
	MUL R9, R6
	MOVD C4, R10
	ADD R10, R6

	ADD $0x10, R0
	SUB $0x10, R1
	B loop16

tail_only:
	MOVD R5, ret+24(FP)
	MOVD R6, ret1+32(FP)
	RET
