// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_findOne()
		{
		mock = Mock(ScintillaControl)
		mock.When.SearchText().Return('This is a test.')
		mock.When.GetSelect().Return(#(cpMax: 8, cpMin: 8))
		mock.When.SetVisibleSelect([anyArgs:]).Return(true)
		mock.ScintillaControl_findreplace_options = Record(find: 's')

		mock.When.numOfMatch([anyArgs:]).Return(Object(num: 3, count: 8))
		mock.Eval(ScintillaControl.ScintillaControl_findAndMatch)
		mock.Verify.SetVisibleSelect(12, 1)

		mock.When.numOfMatch([anyArgs:]).Return(Object(num: 2, count: 8))
		mock.Eval(ScintillaControl.ScintillaControl_findAndMatch, true)
		mock.Verify.SetVisibleSelect(6, 1)
		}
	}