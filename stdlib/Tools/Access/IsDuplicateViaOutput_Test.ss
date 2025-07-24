// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(num) key(num)', [num: 1001])
		Assert(IsDuplicateViaOutput(table, 'num', 1000) is: false)
		QueryOutput(table, [num: 1000])
		Assert(IsDuplicateViaOutput(table, 'num', 1000))

		query = table $ ' where num < 1000'
		Assert(IsDuplicateViaQuery(query, 'num', 1000) is: false) // query fails
		Assert(IsDuplicateViaOutput(query, 'num', 1000))
		}
	}