// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_MakeSummary()
		{
		mock = Mock()
		mock.When.Get().Return('')
		Assert(mock.Eval(ScintillaAddonsControl.MakeSummary) is: '')

		mock.When.Get().Return('hello')
		Assert(mock.Eval(ScintillaAddonsControl.MakeSummary) is: 'hello')

		mock.When.Get().Return(' hello \r\nworld')
		Assert(mock.Eval(ScintillaAddonsControl.MakeSummary) is: 'hello...')

		mock.When.Get().Return('hello'.Repeat(20))
		Assert(mock.Eval(ScintillaAddonsControl.MakeSummary)
			is: 'hello'.Repeat(12) $ '...')
		}
	}