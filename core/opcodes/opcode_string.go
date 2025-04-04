// Code generated by "stringer -type=Opcode"; DO NOT EDIT.

package opcodes

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nop-0]
	_ = x[Pop-1]
	_ = x[Int-2]
	_ = x[Value-3]
	_ = x[True-4]
	_ = x[False-5]
	_ = x[Zero-6]
	_ = x[One-7]
	_ = x[MinusOne-8]
	_ = x[MaxInt-9]
	_ = x[EmptyStr-10]
	_ = x[PushReturn-11]
	_ = x[Load-12]
	_ = x[Store-13]
	_ = x[LoadStore-14]
	_ = x[Dyload-15]
	_ = x[Global-16]
	_ = x[Get-17]
	_ = x[Put-18]
	_ = x[GetPut-19]
	_ = x[RangeTo-20]
	_ = x[RangeLen-21]
	_ = x[This-22]
	_ = x[Is-23]
	_ = x[Isnt-24]
	_ = x[Match-25]
	_ = x[MatchNot-26]
	_ = x[Lt-27]
	_ = x[Lte-28]
	_ = x[Gt-29]
	_ = x[Gte-30]
	_ = x[Add-31]
	_ = x[Sub-32]
	_ = x[Cat-33]
	_ = x[Mul-34]
	_ = x[Div-35]
	_ = x[Mod-36]
	_ = x[LeftShift-37]
	_ = x[RightShift-38]
	_ = x[BitOr-39]
	_ = x[BitAnd-40]
	_ = x[BitXor-41]
	_ = x[BitNot-42]
	_ = x[Not-43]
	_ = x[UnaryPlus-44]
	_ = x[UnaryMinus-45]
	_ = x[InRange-46]
	_ = x[Or-47]
	_ = x[And-48]
	_ = x[Bool-49]
	_ = x[QMark-50]
	_ = x[In-51]
	_ = x[Cover-52]
	_ = x[Jump-53]
	_ = x[JumpTrue-54]
	_ = x[JumpFalse-55]
	_ = x[JumpIs-56]
	_ = x[JumpIsnt-57]
	_ = x[JumpLt-58]
	_ = x[Iter-59]
	_ = x[Iter2-60]
	_ = x[ForIn-61]
	_ = x[ForIn2-62]
	_ = x[ForRange-63]
	_ = x[ForRangeVar-64]
	_ = x[Throw-65]
	_ = x[Try-66]
	_ = x[Catch-67]
	_ = x[CallFuncDiscard-68]
	_ = x[CallFuncNoNil-69]
	_ = x[CallFuncNilOk-70]
	_ = x[CallMethDiscard-71]
	_ = x[CallMethNoNil-72]
	_ = x[CallMethNilOk-73]
	_ = x[Super-74]
	_ = x[Return-75]
	_ = x[ReturnNil-76]
	_ = x[ReturnThrow-77]
	_ = x[ReturnMulti-78]
	_ = x[Closure-79]
	_ = x[BlockBreak-80]
	_ = x[BlockContinue-81]
	_ = x[BlockReturn-82]
	_ = x[BlockReturnNil-83]
}

const _Opcode_name = "NopPopIntValueTrueFalseZeroOneMinusOneMaxIntEmptyStrPushReturnLoadStoreLoadStoreDyloadGlobalGetPutGetPutRangeToRangeLenThisIsIsntMatchMatchNotLtLteGtGteAddSubCatMulDivModLeftShiftRightShiftBitOrBitAndBitXorBitNotNotUnaryPlusUnaryMinusInRangeOrAndBoolQMarkInCoverJumpJumpTrueJumpFalseJumpIsJumpIsntJumpLtIterIter2ForInForIn2ForRangeForRangeVarThrowTryCatchCallFuncDiscardCallFuncNoNilCallFuncNilOkCallMethDiscardCallMethNoNilCallMethNilOkSuperReturnReturnNilReturnThrowReturnMultiClosureBlockBreakBlockContinueBlockReturnBlockReturnNil"

var _Opcode_index = [...]uint16{0, 3, 6, 9, 14, 18, 23, 27, 30, 38, 44, 52, 62, 66, 71, 80, 86, 92, 95, 98, 104, 111, 119, 123, 125, 129, 134, 142, 144, 147, 149, 152, 155, 158, 161, 164, 167, 170, 179, 189, 194, 200, 206, 212, 215, 224, 234, 241, 243, 246, 250, 255, 257, 262, 266, 274, 283, 289, 297, 303, 307, 312, 317, 323, 331, 342, 347, 350, 355, 370, 383, 396, 411, 424, 437, 442, 448, 457, 468, 479, 486, 496, 509, 520, 534}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
