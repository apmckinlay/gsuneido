// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typealgebra

func ResultTypeOfOp(op string) DynType {
	switch op {
	case "Cat", "CatEq":
		return TString
	case "Add", "Sub", "Mul", "Div", "Mod",
		"AddEq", "SubEq", "MulEq", "DivEq", "ModEq",
		"BitAnd", "BitOr", "BitXor",
		"BitAndEq", "BitOrEq", "BitXorEq",
		"LShift", "RShift",
		"LShiftEq", "RShiftEq",
		"Inc", "Dec",
		"PostInc", "PostDec",
		"BitNot":
		return TNumber
	case "And", "Or", "Not",
		"Gt", "Lt", "Gte", "Lte",
		"Is", "Isnt",
		"Match", "MatchNot":
		return TBoolean
	}
	return TUnknown
}
