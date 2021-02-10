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
	_ = x[Whitespace-6]
	_ = x[Comment-7]
	_ = x[Newline-8]
	_ = x[Hash-9]
	_ = x[Comma-10]
	_ = x[Semicolon-11]
	_ = x[At-12]
	_ = x[LParen-13]
	_ = x[RParen-14]
	_ = x[LBracket-15]
	_ = x[RBracket-16]
	_ = x[LCurly-17]
	_ = x[RCurly-18]
	_ = x[RangeTo-19]
	_ = x[RangeLen-20]
	_ = x[OpsStart-21]
	_ = x[Not-22]
	_ = x[BitNot-23]
	_ = x[New-24]
	_ = x[Dot-25]
	_ = x[CompareStart-26]
	_ = x[Is-27]
	_ = x[Isnt-28]
	_ = x[Match-29]
	_ = x[MatchNot-30]
	_ = x[Lt-31]
	_ = x[Lte-32]
	_ = x[Gt-33]
	_ = x[Gte-34]
	_ = x[CompareEnd-35]
	_ = x[QMark-36]
	_ = x[Colon-37]
	_ = x[AssocStart-38]
	_ = x[And-39]
	_ = x[Or-40]
	_ = x[BitOr-41]
	_ = x[BitAnd-42]
	_ = x[BitXor-43]
	_ = x[Add-44]
	_ = x[Sub-45]
	_ = x[Cat-46]
	_ = x[Mul-47]
	_ = x[Div-48]
	_ = x[AssocEnd-49]
	_ = x[Mod-50]
	_ = x[LShift-51]
	_ = x[RShift-52]
	_ = x[Inc-53]
	_ = x[PostInc-54]
	_ = x[Dec-55]
	_ = x[PostDec-56]
	_ = x[AssignStart-57]
	_ = x[Eq-58]
	_ = x[AddEq-59]
	_ = x[SubEq-60]
	_ = x[CatEq-61]
	_ = x[MulEq-62]
	_ = x[DivEq-63]
	_ = x[ModEq-64]
	_ = x[LShiftEq-65]
	_ = x[RShiftEq-66]
	_ = x[BitOrEq-67]
	_ = x[BitAndEq-68]
	_ = x[BitXorEq-69]
	_ = x[AssignEnd-70]
	_ = x[In-71]
	_ = x[Break-72]
	_ = x[Case-73]
	_ = x[Catch-74]
	_ = x[Class-75]
	_ = x[Continue-76]
	_ = x[Default-77]
	_ = x[Do-78]
	_ = x[Else-79]
	_ = x[False-80]
	_ = x[For-81]
	_ = x[Forever-82]
	_ = x[Function-83]
	_ = x[If-84]
	_ = x[Return-85]
	_ = x[Switch-86]
	_ = x[Super-87]
	_ = x[This-88]
	_ = x[Throw-89]
	_ = x[True-90]
	_ = x[Try-91]
	_ = x[While-92]
	_ = x[QueryStart-93]
	_ = x[SummarizeStart-94]
	_ = x[Average-95]
	_ = x[Count-96]
	_ = x[List-97]
	_ = x[Max-98]
	_ = x[Min-99]
	_ = x[Total-100]
	_ = x[SummarizeEnd-101]
	_ = x[Alter-102]
	_ = x[By-103]
	_ = x[Cascade-104]
	_ = x[Create-105]
	_ = x[Delete-106]
	_ = x[Drop-107]
	_ = x[Ensure-108]
	_ = x[Extend-109]
	_ = x[History-110]
	_ = x[Index-111]
	_ = x[Insert-112]
	_ = x[Intersect-113]
	_ = x[Into-114]
	_ = x[Join-115]
	_ = x[Key-116]
	_ = x[Leftjoin-117]
	_ = x[Lower-118]
	_ = x[Minus-119]
	_ = x[Project-120]
	_ = x[Remove-121]
	_ = x[Rename-122]
	_ = x[Reverse-123]
	_ = x[Set-124]
	_ = x[Sort-125]
	_ = x[Summarize-126]
	_ = x[Sview-127]
	_ = x[Times-128]
	_ = x[To-129]
	_ = x[Union-130]
	_ = x[Unique-131]
	_ = x[Update-132]
	_ = x[View-133]
	_ = x[Where-134]
	_ = x[Ntokens-135]
}

const _Token_name = "NilEofErrorIdentifierNumberStringWhitespaceCommentNewlineHashCommaSemicolonAtLParenRParenLBracketRBracketLCurlyRCurlyRangeToRangeLenOpsStartNotBitNotNewDotCompareStartIsIsntMatchMatchNotLtLteGtGteCompareEndQMarkColonAssocStartAndOrBitOrBitAndBitXorAddSubCatMulDivAssocEndModLShiftRShiftIncPostIncDecPostDecAssignStartEqAddEqSubEqCatEqMulEqDivEqModEqLShiftEqRShiftEqBitOrEqBitAndEqBitXorEqAssignEndInBreakCaseCatchClassContinueDefaultDoElseFalseForForeverFunctionIfReturnSwitchSuperThisThrowTrueTryWhileQueryStartSummarizeStartAverageCountListMaxMinTotalSummarizeEndAlterByCascadeCreateDeleteDropEnsureExtendHistoryIndexInsertIntersectIntoJoinKeyLeftjoinLowerMinusProjectRemoveRenameReverseSetSortSummarizeSviewTimesToUnionUniqueUpdateViewWhereNtokens"

var _Token_index = [...]uint16{0, 3, 6, 11, 21, 27, 33, 43, 50, 57, 61, 66, 75, 77, 83, 89, 97, 105, 111, 117, 124, 132, 140, 143, 149, 152, 155, 167, 169, 173, 178, 186, 188, 191, 193, 196, 206, 211, 216, 226, 229, 231, 236, 242, 248, 251, 254, 257, 260, 263, 271, 274, 280, 286, 289, 296, 299, 306, 317, 319, 324, 329, 334, 339, 344, 349, 357, 365, 372, 380, 388, 397, 399, 404, 408, 413, 418, 426, 433, 435, 439, 444, 447, 454, 462, 464, 470, 476, 481, 485, 490, 494, 497, 502, 512, 526, 533, 538, 542, 545, 548, 553, 565, 570, 572, 579, 585, 591, 595, 601, 607, 614, 619, 625, 634, 638, 642, 645, 653, 658, 663, 670, 676, 682, 689, 692, 696, 705, 710, 715, 717, 722, 728, 734, 738, 743, 750}

func (i Token) String() string {
	if i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
