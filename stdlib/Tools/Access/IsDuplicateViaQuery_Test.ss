// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(num) key(num)', [num: 1001])
		Assert(IsDuplicateViaQuery(table, 'num', 1000) is: false)
		Assert(IsDuplicateViaQuery(table $ ' sort num', 'num', 1000) is: false)
		QueryOutput(table, [num: 1000])
		Assert(IsDuplicateViaQuery(table, 'num', 1000))
		Assert(IsDuplicateViaQuery(table $ ' sort num', 'num', 1000))
		}
	}