// Code generated by "stringer -type=Token"; DO NOT EDIT.

package tokens

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nil-0]
	_ = x[Eof-1]
	_ = x[Error-2]
	_ = x[Identifier-3]
	_ = x[Number-4]
	_ = x[String-5]
	_ = x[Symbol-6]
	_ = x[Whitespace-7]
	_ = x[Comment-8]
	_ = x[Newline-9]
	_ = x[Hash-10]
	_ = x[Comma-11]
	_ = x[Semicolon-12]
	_ = x[At-13]
	_ = x[LParen-14]
	_ = x[RParen-15]
	_ = x[LBracket-16]
	_ = x[RBracket-17]
	_ = x[LCurly-18]
	_ = x[RCurly-19]
	_ = x[RangeTo-20]
	_ = x[RangeLen-21]
	_ = x[OpsStart-22]
	_ = x[Not-23]
	_ = x[BitNot-24]
	_ = x[New-25]
	_ = x[Dot-26]
	_ = x[CompareStart-27]
	_ = x[Is-28]
	_ = x[Isnt-29]
	_ = x[Match-30]
	_ = x[MatchNot-31]
	_ = x[Lt-32]
	_ = x[Lte-33]
	_ = x[Gt-34]
	_ = x[Gte-35]
	_ = x[CompareEnd-36]
	_ = x[QMark-37]
	_ = x[Colon-38]
	_ = x[AssocStart-39]
	_ = x[And-40]
	_ = x[Or-41]
	_ = x[BitOr-42]
	_ = x[BitAnd-43]
	_ = x[BitXor-44]
	_ = x[Add-45]
	_ = x[Sub-46]
	_ = x[Cat-47]
	_ = x[Mul-48]
	_ = x[Div-49]
	_ = x[AssocEnd-50]
	_ = x[Mod-51]
	_ = x[LShift-52]
	_ = x[RShift-53]
	_ = x[Inc-54]
	_ = x[PostInc-55]
	_ = x[Dec-56]
	_ = x[PostDec-57]
	_ = x[AssignStart-58]
	_ = x[Eq-59]
	_ = x[AddEq-60]
	_ = x[SubEq-61]
	_ = x[CatEq-62]
	_ = x[MulEq-63]
	_ = x[DivEq-64]
	_ = x[ModEq-65]
	_ = x[LShiftEq-66]
	_ = x[RShiftEq-67]
	_ = x[BitOrEq-68]
	_ = x[BitAndEq-69]
	_ = x[BitXorEq-70]
	_ = x[AssignEnd-71]
	_ = x[In-72]
	_ = x[Break-73]
	_ = x[Case-74]
	_ = x[Catch-75]
	_ = x[Class-76]
	_ = x[Continue-77]
	_ = x[Default-78]
	_ = x[Do-79]
	_ = x[Else-80]
	_ = x[False-81]
	_ = x[For-82]
	_ = x[Forever-83]
	_ = x[Function-84]
	_ = x[If-85]
	_ = x[Return-86]
	_ = x[Switch-87]
	_ = x[Super-88]
	_ = x[This-89]
	_ = x[Throw-90]
	_ = x[True-91]
	_ = x[Try-92]
	_ = x[While-93]
	_ = x[QueryStart-94]
	_ = x[SummarizeStart-95]
	_ = x[Average-96]
	_ = x[Count-97]
	_ = x[List-98]
	_ = x[Max-99]
	_ = x[Min-100]
	_ = x[Total-101]
	_ = x[SummarizeEnd-102]
	_ = x[Alter-103]
	_ = x[By-104]
	_ = x[Cascade-105]
	_ = x[Create-106]
	_ = x[Delete-107]
	_ = x[Drop-108]
	_ = x[Ensure-109]
	_ = x[Extend-110]
	_ = x[History-111]
	_ = x[Index-112]
	_ = x[Insert-113]
	_ = x[Intersect-114]
	_ = x[Into-115]
	_ = x[Join-116]
	_ = x[Key-117]
	_ = x[Leftjoin-118]
	_ = x[Lower-119]
	_ = x[Minus-120]
	_ = x[Project-121]
	_ = x[Remove-122]
	_ = x[Rename-123]
	_ = x[Reverse-124]
	_ = x[Set-125]
	_ = x[Sort-126]
	_ = x[Summarize-127]
	_ = x[Sview-128]
	_ = x[Times-129]
	_ = x[To-130]
	_ = x[Union-131]
	_ = x[Unique-132]
	_ = x[Update-133]
	_ = x[View-134]
	_ = x[Where-135]
	_ = x[Ntokens-136]
}

const _Token_name = "NilEofErrorIdentifierNumberStringSymbolWhitespaceCommentNewlineHashCommaSemicolonAtLParenRParenLBracketRBracketLCurlyRCurlyRangeToRangeLenOpsStartNotBitNotNewDotCompareStartIsIsntMatchMatchNotLtLteGtGteCompareEndQMarkColonAssocStartAndOrBitOrBitAndBitXorAddSubCatMulDivAssocEndModLShiftRShiftIncPostIncDecPostDecAssignStartEqAddEqSubEqCatEqMulEqDivEqModEqLShiftEqRShiftEqBitOrEqBitAndEqBitXorEqAssignEndInBreakCaseCatchClassContinueDefaultDoElseFalseForForeverFunctionIfReturnSwitchSuperThisThrowTrueTryWhileQueryStartSummarizeStartAverageCountListMaxMinTotalSummarizeEndAlterByCascadeCreateDeleteDropEnsureExtendHistoryIndexInsertIntersectIntoJoinKeyLeftjoinLowerMinusProjectRemoveRenameReverseSetSortSummarizeSviewTimesToUnionUniqueUpdateViewWhereNtokens"

var _Token_index = [...]uint16{0, 3, 6, 11, 21, 27, 33, 39, 49, 56, 63, 67, 72, 81, 83, 89, 95, 103, 111, 117, 123, 130, 138, 146, 149, 155, 158, 161, 173, 175, 179, 184, 192, 194, 197, 199, 202, 212, 217, 222, 232, 235, 237, 242, 248, 254, 257, 260, 263, 266, 269, 277, 280, 286, 292, 295, 302, 305, 312, 323, 325, 330, 335, 340, 345, 350, 355, 363, 371, 378, 386, 394, 403, 405, 410, 414, 419, 424, 432, 439, 441, 445, 450, 453, 460, 468, 470, 476, 482, 487, 491, 496, 500, 503, 508, 518, 532, 539, 544, 548, 551, 554, 559, 571, 576, 578, 585, 591, 597, 601, 607, 613, 620, 625, 631, 640, 644, 648, 651, 659, 664, 669, 676, 682, 688, 695, 698, 702, 711, 716, 721, 723, 728, 734, 740, 744, 749, 756}

func (i Token) String() string {
	if i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
