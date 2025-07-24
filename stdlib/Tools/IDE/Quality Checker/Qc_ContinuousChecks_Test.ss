// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_skipQc?()
		{
		skip = Qc_ContinuousChecks.Qc_ContinuousChecks_skipQc?
		Assert(skip(0, 0, ''))
		Assert(skip('', 0, ''))
		Assert(skip('', '', ''))
		code = 'class
			{
			a = 10
			b = 20
			}'
		Assert(skip('a', 'b', code) is: false)
		Assert(skip('', '', '{}'))
		Assert(skip('a', 'b.js', code))
		Assert(skip('a', 'b.css', code))
		Assert(skip('a', 'b', '#(
			a: 10
			b: 20
			)'))
		}
	}