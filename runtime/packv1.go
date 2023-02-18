// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !neworder

package runtime

// Packed values start with one of the following type tags,
// except for the special case of a zero length string
// which is encoded as a zero length buffer.
// NOTE: this order is significant, it determines sorting
const (
	PackFalse = iota
	PackTrue
	PackMinus
	PackPlus
	PackString
	PackDate
	PackObject
	PackRecord
)

// all new values so we can interact between old & new
const (
	PackFalseOther  = packBool << packShift
	PackTrueOther   = packBool<<packShift | 1
	PackMinusOther  = packNum << packShift
	PackPlusOther   = packNum<<packShift | 1
	PackStringOther = packStr << packShift
	PackDateOther   = packDate << packShift
	PackObjectOther = packObject << packShift
	PackRecordOther = packObject<<packShift | 1
)

const (
	packBool = iota + 1
	packNum  // SuInt, SuDnum
	packStr  // SuStr, SuConcat, SuExcept
	packDate
	packObject
)
const packShift = 4

func PackedOrd(s string) Ord {
	if s == "" {
		return ordStr
	}
	switch s[0] {
	case PackFalse:
		return ordBool
	case PackTrue:
		return ordBool
	case PackMinus:
		return ordNum
	case PackPlus:
		return ordNum
	case PackString:
		return ordStr
	case PackDate:
		return ordDate
	case PackObject:
		return ordObject
	case PackRecord:
		return ordObject
	}
	panic("unknown")
}
