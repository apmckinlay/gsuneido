// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// ensure cursors don't use tempindexes
Test
	{
	Test_main()
		{
		query = "stdlib where name > 'X' sort num"
		Assert(QueryStrategy(query) has: "tempindex")
		Cursor(query)
			{|c|
			Assert(QueryStrategy(c) hasnt: "tempindex")
			}
		}
	}