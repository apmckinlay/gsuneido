// Code generated by "stringer -type=Opcode"; DO NOT EDIT.

package opcodes

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nop-0]
	_ = x[Pop-1]
	_ = x[Dup-2]
	_ = x[Swap-3]
	_ = x[Int-4]
	_ = x[Value-5]
	_ = x[True-6]
	_ = x[False-7]
	_ = x[Zero-8]
	_ = x[One-9]
	_ = x[MinusOne-10]
	_ = x[MaxInt-11]
	_ = x[EmptyStr-12]
	_ = x[Load-13]
	_ = x[Store-14]
	_ = x[LoadStore-15]
	_ = x[Dyload-16]
	_ = x[Global-17]
	_ = x[Get-18]
	_ = x[Put-19]
	_ = x[GetPut-20]
	_ = x[RangeTo-21]
	_ = x[RangeLen-22]
	_ = x[This-23]
	_ = x[Is-24]
	_ = x[Isnt-25]
	_ = x[Match-26]
	_ = x[MatchNot-27]
	_ = x[Lt-28]
	_ = x[Lte-29]
	_ = x[Gt-30]
	_ = x[Gte-31]
	_ = x[Add-32]
	_ = x[Sub-33]
	_ = x[Cat-34]
	_ = x[Mul-35]
	_ = x[Div-36]
	_ = x[Mod-37]
	_ = x[LeftShift-38]
	_ = x[RightShift-39]
	_ = x[BitOr-40]
	_ = x[BitAnd-41]
	_ = x[BitXor-42]
	_ = x[BitNot-43]
	_ = x[Not-44]
	_ = x[UnaryPlus-45]
	_ = x[UnaryMinus-46]
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
	_ = x[Iter-58]
	_ = x[ForIn-59]
	_ = x[Throw-60]
	_ = x[Try-61]
	_ = x[Catch-62]
	_ = x[CallFuncDiscard-63]
	_ = x[CallFuncNoNil-64]
	_ = x[CallFuncNilOk-65]
	_ = x[CallMethDiscard-66]
	_ = x[CallMethNoNil-67]
	_ = x[CallMethNilOk-68]
	_ = x[Super-69]
	_ = x[Return-70]
	_ = x[ReturnNil-71]
	_ = x[Closure-72]
	_ = x[BlockBreak-73]
	_ = x[BlockContinue-74]
	_ = x[BlockReturn-75]
	_ = x[BlockReturnNil-76]
}

const _Opcode_name = "NopPopDupSwapIntValueTrueFalseZeroOneMinusOneMaxIntEmptyStrLoadStoreLoadStoreDyloadGlobalGetPutGetPutRangeToRangeLenThisIsIsntMatchMatchNotLtLteGtGteAddSubCatMulDivModLeftShiftRightShiftBitOrBitAndBitXorBitNotNotUnaryPlusUnaryMinusOrAndBoolQMarkInCoverJumpJumpTrueJumpFalseJumpIsJumpIsntIterForInThrowTryCatchCallFuncDiscardCallFuncNoNilCallFuncNilOkCallMethDiscardCallMethNoNilCallMethNilOkSuperReturnReturnNilClosureBlockBreakBlockContinueBlockReturnBlockReturnNil"

var _Opcode_index = [...]uint16{0, 3, 6, 9, 13, 16, 21, 25, 30, 34, 37, 45, 51, 59, 63, 68, 77, 83, 89, 92, 95, 101, 108, 116, 120, 122, 126, 131, 139, 141, 144, 146, 149, 152, 155, 158, 161, 164, 167, 176, 186, 191, 197, 203, 209, 212, 221, 231, 233, 236, 240, 245, 247, 252, 256, 264, 273, 279, 287, 291, 296, 301, 304, 309, 324, 337, 350, 365, 378, 391, 396, 402, 411, 418, 428, 441, 452, 466}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
