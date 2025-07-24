// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		tbl = .MakeTable("(a,b) key(a)")
		Assert(QueryAll(tbl) is: [])
		Assert(QueryAll(tbl, a: 1) is: [])

		QueryOutput(tbl, [a: 1, b: 2])
		Assert(QueryAll(tbl, a: 9) is: [])
		Assert(QueryAll(tbl) is: [[a: 1, b: 2]])
		Assert(QueryAll(tbl, a: 1) is: [[a: 1, b: 2]])

		QueryOutput(tbl, [a: 2, b: 2])
		Assert(QueryAll(tbl, a: 9) is: [])
		Assert(QueryAll(tbl, a: 1) is: [[a: 1, b: 2]])
		Assert(QueryAll(tbl) is: [[a: 1, b: 2], [a: 2, b: 2]])
		Assert(QueryAll(tbl, b: 2) is: [[a: 1, b: 2], [a: 2, b: 2]])

		// test limit of 4, so last rec will not be included
		QueryOutput(tbl, [a: 3, b: 2])
		QueryOutput(tbl, [a: 4, b: 2])
		QueryOutput(tbl, [a: 5, b: 2])

		Assert(QueryAll(tbl, limit: 4)
			is: [[a: 1, b: 2], [a: 2, b: 2], [a: 3, b: 2], [a: 4, b: 2]])
		}
	}
