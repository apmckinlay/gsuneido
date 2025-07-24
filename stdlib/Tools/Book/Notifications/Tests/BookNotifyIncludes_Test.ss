// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(not BookNotifyIncludes?("", "xxx"))
		Assert(not BookNotifyIncludes?(#(), "xxx"))
		Assert(BookNotifyIncludes?("abc", "abc"))
		Assert(BookNotifyIncludes?(#(abc), "abc"))
		Assert(BookNotifyIncludes?("ab,de,fg", "de"))
		Assert(BookNotifyIncludes?(#(ab,de,fg), "de"))
		Assert(not BookNotifyIncludes?("ab,de,fg", "d"))
		Assert(not BookNotifyIncludes?(#(ab,de,fg), "e"))
		}
	}