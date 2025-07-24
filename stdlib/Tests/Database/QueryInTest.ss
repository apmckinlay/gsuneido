// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		tbl = .MakeTable('(a) key(a)', [a: 1], [a: 2], [a: 3], [a: 4])
		Assert(QueryList(tbl, #a) equalsSet: [1, 2, 3, 4])
		Assert(QueryList(tbl $ ' where a in (2,3)', #a) equalsSet: [2, 3])
		Assert(QueryList(tbl $ ' where a not in (2,3)', #a) equalsSet: [1, 4])
		}
	}
