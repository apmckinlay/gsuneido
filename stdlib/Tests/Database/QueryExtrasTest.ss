// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		tbl = .MakeTable("(a) key(a)")
		for a in ..10
			QueryOutput(tbl, [:a])
		Assert(QueryMin(tbl, 'a') is: 0)
		Assert(QueryMax(tbl, 'a', '') is: 9)
		Assert(QueryCount(tbl $ ' where a < 5') is: 5)
		Assert(QueryTotal(tbl $ ' where a < 5', 'a') is: 10)
		Assert(QueryList(tbl $ ' where a < 4', 'a') equalsSet: #(0,1,2,3))
		Transaction(read:)
			{ |t|
			Assert(t.QueryMin(tbl, 'a') is: 0)
			Assert(t.QueryMax(tbl, 'a') is: 9)
			Assert(t.QueryCount(tbl $ ' where a < 5') is: 5)
			Assert(t.QueryTotal(tbl $ ' where a <= 3', 'a') is: 6)
			Assert(t.QueryList(tbl $ ' where a < 4', 'a') equalsSet: #(0,1,2,3))
			}
		}
	}