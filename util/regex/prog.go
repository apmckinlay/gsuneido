// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate stringer -type=opType

type Pattern string

type opType byte

const (
	_                opType = iota
	opChar                  // char
	opCharIgnoreCase        // char
	opJump                  // int16
	opSplitNext             // int16
	opSplitJump             // int16
	opAnyNotNL              //
	opHalfSet               // [16]byte
	opFullSet               // [32]byte
	opListSet               // uint8 []byte
	opWordStart             //
	opWordEnd               //
	opLineStart             //
	opLineEnd               //
	opStrStart              //
	opStrEnd                //
	opSave                  // byte
	opDoneSave1             //
	opOnePass               //
	opBranchNext            // int16
	opBranchJump            // int16
	opLiteralSubstr         // []byte (to end)
	opLiteralPrefix         // []byte (to end)
	opLiteralSuffix         // []byte (to end)
	opLiteralEqual          // []byte (to end)
	opPrefix                // uint8 []byte
)

func (pat Pattern) String() string {
	var sb strings.Builder
	pi := 0
	for pi < len(pat) {
		inc, s := pat.opstr(pi)
		sb.WriteString(fmt.Sprintf("%d: %s\n", pi, s))
		pi += inc
	}
	return sb.String()
}

func (pat Pattern) opstr(pi int) (int, string) {
	op := opType(pat[pi])
	opstr := op.String()[2:]
	switch op {
	case opChar, opCharIgnoreCase:
		return 2, fmt.Sprintf("Char %c", pat[pi+1])
	case opJump, opSplitJump, opSplitNext, opBranchNext, opBranchJump:
		jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
		return 3, fmt.Sprint(opstr, " ", pi+int(jmp))
	case opHalfSet:
		return 1 + 16, opstr
	case opFullSet:
		return 1 + 32, opstr
	case opSave:
		return 2, fmt.Sprintf("Save %d", int(pat[pi+1]))
	case opListSet, opPrefix:
		n := int(pat[pi+1])
		return n + 2, fmt.Sprintf("%s %q", opstr, string(pat[pi+2:pi+2+n]))
	case opLiteralSubstr, opLiteralPrefix, opLiteralSuffix, opLiteralEqual:
		return len(pat), fmt.Sprintf("%s %q", opstr, string(pat[pi+1:]))
	default:
		return 1, opstr
	}
}

func (pat Pattern) opstr1(pi int16) string {
	_, s := pat.opstr(int(pi))
	return strconv.Itoa(int(pi)) + ": " + s
}
