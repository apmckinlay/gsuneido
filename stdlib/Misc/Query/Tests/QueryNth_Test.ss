// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		tbl = .MakeTable("(a) key(a)")
		for (i = 0; i < 10; ++i)
			QueryOutput(tbl, [a: i])
		query = tbl $ ' sort a'
		Assert(QueryNth(0, query).a is: 0)
		Assert(QueryNth(3, query).a is: 3)
		Assert(QueryNth(999999, query) is: false)
		}
	}
