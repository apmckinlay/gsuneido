// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_askReplaceSection()
		{
		mock = MockObject(#(
			(SetFocus),
			(ReplaceSel, true),
			(SetFocus)))
		FormulaEditor.FormulaEditor_askReplaceSection(mock, { true })

		mock = MockObject(#(
			(SetFocus)))
		FormulaEditor.FormulaEditor_askReplaceSection(mock, { false })
		}
	}
