// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Detab()
		{
		fn = Md_Helper.Detab
		Assert(fn('') is: '')
		Assert(fn('foo  \tbar\t') is: 'foo  \tbar\t')
		Assert(fn('  foo  \tbar\t') is: '  foo  \tbar\t')
		Assert(fn('  \tfoo  \tbar\t') is: '    foo  \tbar\t')
		Assert(fn('  \t \tfoo  \tbar\t') is: '        foo  \tbar\t')
		}
	}