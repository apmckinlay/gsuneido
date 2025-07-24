// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_findMatch()
		{
		fn = AutoFillFieldControl.AutoFillFieldControl_findMatch

		Assert(fn('a', '', false) is: false)
		Assert(fn('a', '', #(a, b)) is: false)
		Assert(fn('c', '', #(a, b)) is: false)
		Assert(fn('', '', #(aa, bb)) is: false)

		Assert(fn('a', '', #(aa, bb)) is: 'aa')
		Assert(fn('A', '', #(aa, bb)) is: 'aa')

		Assert(fn('a', 'a', #(aa, bb)) is: false)
		// current value is a prefix of the prev value, which means the user is deleting
		Assert(fn('a', 'aa', #(aa, bb)) is: false)
		Assert(fn('a', 'bb', #(aaa, bbb)) is: 'aaa')
		Assert(fn('aa', 'a', #(aaa, bbb)) is: 'aaa')
		}
	}