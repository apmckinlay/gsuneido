// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_hasSimilarWord?()
		{
		specialChars = Diff2Model.Diff2Model_specialChar
		.check([], [], false)
		.check(specialChars, specialChars, false)
		.check(['abc'], ['abc'], true)
		.check(['abc', '.'], ['abc'], true)
		.check(['abc', 'def', '\t'].MergeUnion(specialChars),
			['ghi', 'jkl', '\t'].MergeUnion(specialChars), false)
		}

	check(tokens1, tokens2, expected)
		{
		fn = Diff2Model.Diff2Model_hasSimilarWord?
		Assert(fn(Diff.SideBySide(tokens1, tokens2)) is: expected)
		}

	Test_getTokens()
		{
		fn = Diff2Model.Diff2Model_getTokens

		Assert(fn("") is: #(#(), addr: #(0)))
		Assert(fn("tokens test") is: #(#("tokens", " ", "test"), addr: #(0, 6, 7, 11)))
		Assert(fn("./:;+=") is: #(#(".", "/", ":", ";", "+="), addr: #(0, 1, 2, 3, 4, 6)))
		Assert(fn("Dot.(Brackets)") is: #(#("Dot", ".", "(", "Brackets", ")"),
			addr: #(0, 3, 4, 5, 13, 14)))
		}

	Test_GetRowIndics()
		{
		// NOTE - Diff algorithm returns differences from end of line to beginning

		// general word changes and <= 40% of text words changed
		.checkIndics("this isnt your ext lines", "this isnt your text line",
			[[19, 5], [15, 3]], [[20, 4], [15, 4]])
		// more than 40% of text words changed
		.checkIndics("a test line here", "two test lines here", [], [])
		// character specific highlight (only one word changed/consecutive changes)
		.checkIndics("this is a test", "this isnt a test", [], [[8, 1], [7, 1]])
		// non-consecutive character-specific changes = full word highlight
		.checkIndics("result test", "TheResultRow test", [[0, 6]], [[0, 12]])
		// one word changed, but only one word on line = no highlighting
		.checkIndics("hello", "}", [], [])
		// exactly one deletion
		.checkIndics("this is a test", "this is a ", [[10, 4]], [])
		// exactly one insertion
		.checkIndics("this is a test", "when this is a test", [], [[4, 1], [0, 4]])
		// whitespace change only
		.checkIndics("test", "		test", [], [[1, 1], [0, 1]])
		// whitespace and character change
		.checkIndics("	hi hello test", "			hi hey test",
			[[4, 5]], [[6, 3], [2, 1], [1, 1]])
		// no similar words, only similar characters
		.checkIndics("this.is(a/test)", "no.similar(words/here)", [], [])
		// too many changes = no highlighting
		.checkIndics("this is a really long line", "the whole line is changed", [], [])
		// numChanges < 1 on one side but > 1 on the other side - treat as > 1
		.checkIndics("similar.test here", "similar.is new here", [[8, 4]],
			[[14, 1], [11, 3], [8, 2]])
		// only special characters, don't highlight
		.checkIndics("./()", "[]{}", [], [])
		}

	Test_consecChrsOnly?()
		{
		fn = Diff2Model.Diff2Model_consecChrsOnly?

		test1 = Diff("testone", "texttwo")
		test2 = Diff("testaaa", "testone")
		test3 = Diff("delete(rec)", "DeleteRow(rec)")
		test4 = Diff("result", "ResultRow")

		Assert(fn(test1) is: false)
		Assert(fn(test2))
		Assert(fn(test3) is: false)
		Assert(fn(test4) is: false)
		}

	checkIndics(textOld, textNew, expectedOld, expectedNew, row = 0)
		{
		textOld = textOld.Lines()
		textNew = textNew.Lines()
		model = new Diff2Model(textOld, textNew)
		rowIndics = model.GetRowIndics(row)
		expectedOld = expectedOld.Map({ it.ListToNamed(#pos, #length) })
		expectedNew = expectedNew.Map({ it.ListToNamed(#pos, #length) })
		Assert(rowIndics.OldToNew is: expectedOld)
		Assert(rowIndics.NewToOld is: expectedNew)
		}
	}