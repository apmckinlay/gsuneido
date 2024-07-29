// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

#include "textflag.h"

// func Read() uint64
TEXT ·Read(SB),NOSPLIT,$0-8
	RDTSC
	SHLQ $32, DX
	ADDQ DX, AX
	MOVQ AX, ret+0(FP)
	RET
