// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build neworder

package runtime

// Packed values start with one of the following type tags,
// except for the special case of a zero length string
// which is encoded as a zero length buffer.
const (
	PackFalse  = packBool << packShift
	PackTrue   = packBool<<packShift | 1
	PackMinus  = packNum << packShift
	PackPlus   = packNum<<packShift | 1
	PackString = packStr << packShift
	PackDate   = packDate << packShift
	PackObject = packObject << packShift
	PackRecord = packObject<<packShift | 1
)

const (
	packBool = iota + 1
	packNum  // SuInt, SuDnum
	packStr  // SuStr, SuConcat, SuExcept
	packDate
	packObject
)
const packShift = 4

const (
	PackFalseOther = iota
	PackTrueOther
	PackMinusOther
	PackPlusOther
	PackStringOther
	PackDateOther
	PackObjectOther
	PackRecordOther
)

func PackedOrd(s string) Ord {
	if s == "" {
		return ordStr
	}
	return Ord(s[0] >> packShift)
}
