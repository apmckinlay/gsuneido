// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a) key(a)',
			[a: 1], [a: 2], [a: 3])
		Assert(QueryFirst(table $ ' sort a') is: [a: 1])
		Assert(QueryLast(table $ ' sort a') is: [a: 3])
		Assert(Query1(table $ ' where a = 2') is: [a: 2])
		Assert(Query1(table $ ' where a = 9') is: false)
		Assert({ Query1(table $ ' where a > 0') } throws: "Query1 not unique")
		Transaction(update:)
			{ |t|
			Assert(t.QueryFirst(table $ ' sort a') is: [a: 1])
			Assert(t.QueryLast(table $ ' sort a') is: [a: 3])
			Assert(t.Query1(table $ ' where a = 2') is: [a: 2])
			Assert(t.Query1(table $ ' where a = 9') is: false)
			Assert({ t.Query1(table $ ' where a > 0') } throws: "Query1 not unique")
			}
		}
	}