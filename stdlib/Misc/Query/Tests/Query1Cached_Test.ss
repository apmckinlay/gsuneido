// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		query = 'stdlib where name = "Init"'
		x = Query1(query)
		Suneido.Delete(#Query1Cache)
		y = Query1Cached(query)
		Assert(x is: y)
		Assert(Suneido.Query1Cache.GetMissRate() is: 1)
		y = Query1Cached(query)
		Assert(x is: y)
		Assert(Suneido.Query1Cache.GetMissRate() is: .5)

		Assert(Query1Cached('stdlib', name: 'Init') is: x)
		}
	}
