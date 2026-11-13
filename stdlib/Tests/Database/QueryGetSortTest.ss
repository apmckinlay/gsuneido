// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Query.GetSort("reverse where reverse = 1 sort reverse one")
			is: "reverse one")
		Assert(Query.GetSort("reverse where reverse = 1 sort \t reverse one, two, three")
			is: "reverse one,two,three")
		Assert(Query.GetSort("table sort reverseField, two")
			is: "reverseField,two")
		Assert(Query.GetSort("table sort table") is: #table)
		Assert(Query.GetSort(#table) is: "")

		Assert(Query.GetSort("foo sort bar /* blah */") is: #bar)
		Assert(Query.GetSort("foo sort bar, baz /* blah*/") is: "bar,baz")
		Assert(Query.GetSort("foo sort bar, /* blah*/ baz /* blah*/") is: "bar,baz")
		Assert(Query.GetSort("foo sort\nbar,\n/* blah */\nbaz\n/* blah*/") is: "bar,baz")
		}
	}
