// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(QueryGetSort('reverse where reverse = 1 sort reverse one')
			is: 'reverse one')
		Assert(QueryGetSort('reverse where reverse = 1 sort \t reverse one, two, three',
			keepReverse:) is: 'reverse one,two,three')
		Assert(QueryGetSort('table sort reverseField, two')
			is: 'reverseField,two')
		Assert(QueryGetSort('table sort table') is: 'table')
		Assert(QueryGetSort('table') is: '')

		Assert(QueryGetSort('foo sort bar /* blah */') is: "bar")
		Assert(QueryGetSort('foo sort bar, baz /* blah*/') is: "bar,baz")
		Assert(QueryGetSort('foo sort bar, /* blah*/ baz /* blah*/') is: "bar,baz")
		Assert(QueryGetSort('foo sort\nbar,\n/* blah */\nbaz\n/* blah*/') is: "bar,baz")
		}
	}