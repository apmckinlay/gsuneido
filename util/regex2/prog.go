// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=opType

type Pattern string

type opType byte

const (
	opChar opType = iota + 1
	opCharIgnoreCase
	opJump
	opSplitFirst
	opSplitLast
	opAny
	opHalfSet
	opFullSet
	opListSet
	opWordStart
	opWordEnd
	opLineStart
	opLineEnd
	opStrStart
	opStrEnd
	opSave
	opStop
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
	case opChar:
		return 2, fmt.Sprintf("Char %c", pat[pi+1])
	case opJump, opSplitFirst, opSplitLast:
		jmp := int16(pat[pi+1])<<8 | int16(pat[pi+2])
		return 3, fmt.Sprint(opstr, " ", pi+int(jmp))
	case opHalfSet:
		return 1+16, opstr
	case opFullSet:
		return 1+32, opstr
	case opListSet:
		n := int(pat[pi+1])
		return n + 2, fmt.Sprintf("List %q", string(pat[pi+2:pi+2+n]))
	case opSave:
		return 2, fmt.Sprintf("Save %d", int(pat[pi+1]))
	default:
		return 1, opstr
	}
}

func (pat Pattern) opstr1(pi int) string {
	_, s := pat.opstr(pi)
	return s
}
