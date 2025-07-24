// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Import1()
		{
		mock = Mock(ImportTab)
		mock.Fields = #('a', 'b')
		Assert(mock.Eval(ImportTab.Import1, '') is: Record())
		Assert(mock.Eval(ImportTab.Import1, 'a\tb') is: Record(a: 'a', b: 'b'))
		Assert(mock.Eval(ImportTab.Import1, '"a"\tb') is: Record(a: 'a', b: 'b'))
		Assert(mock.Eval(ImportTab.Import1, '"a"\tb\tc') is: Record(a: 'a', b: 'b'))
		}
	}