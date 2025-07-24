// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_NamePath()
		{
		.MakeLibraryRecord(
			[num: 1, name: 'one', group: -1],
			[num: 2, name: 'two', group: -1],
			[num: 3, name: 'leaf', group: -1])
		QueryDo('update Test_lib where name is "leaf" set parent = 2')
		QueryDo('update Test_lib where name is "two" set parent = 1')
		QueryDo('update Test_lib where name is "one" set parent = 0')
		Assert(LibHelp.NamePath('Test_lib', 'leaf') is: 'Test_lib/one/two/leaf')
		Assert(LibHelp.NamePath('Test_lib', '_leaf') is: 'Test_lib/one/two/leaf')
		Assert(LibHelp.NamePath('Test_lib', 'two') is: 'Test_lib/one/two')
		Assert(LibHelp.NamePath('Test_lib', 'one') is: 'Test_lib/one')
		}
	}