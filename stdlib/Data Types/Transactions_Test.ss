// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Query1Cached()
		{
		query = 'tables where table = "stdlib"'
		x = Query1(query)
		Transaction(read:)
			{|t|
			y = t.Query1Cached(query)
			Assert(y is: x)
			Assert(.missRate(t) is: 1)
			y = t.Query1Cached(query)
			Assert(y is: x)
			Assert(.missRate(t) is: .5)
			}
		Transaction(read:)
			{|t|
			y = t.Query1Cached(query)
			Assert(.missRate(t) is: 1)
			}
		x = Query1(query)
		Transaction(read:)
			{ |t|
			y = t.Query1Cached(query)
			Assert(y is: x)
			Assert(.missRate(t) is: 1)
			y = t.Query1Cached(query)
			y = t.Query1Cached(query)
			y = t.Query1Cached(query)
			Assert(.missRate(t) is: 0.25)
			}
		}
	missRate(t)
		{
		return t.Data().Query1Cache.GetMissRate()
		}
	Test_QueryAll()
		{
		tbl = .MakeTable('(a) key(a)')
		for a in .. 20
			QueryOutput(tbl, [:a])

		test = {|query, expected|
			Assert(QueryAll(query).Map!({ it.a }) is: expected) // in memory
			Assert(QueryAll(query, 20).Map!({ it.a }) is: expected) // on query
			}
		test(tbl, Seq(20))
		test(tbl $ ' sort a', Seq(20))
		test(tbl $ ' where a < 10 sort reverse a', Seq(10).Reverse!())

		Assert(QueryAll(tbl, 10).Map!({ it.a }) is: Seq(10))
		}
	}
