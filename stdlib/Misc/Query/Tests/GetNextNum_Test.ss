// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Suneido.forceRetryTransaction = true
		gnn = new GetNextNum
			{
			GetNextNum_date()
				{ return .Date }
			}
		gnn.Date = Date()
		gnn.Create('getnextnum_test', nextnum: 1000)
		Assert(Query1('getnextnum_test')
			is: Record(nextnum: 1000, getnextnum_reserved_till: Date.Begin()))
		Assert(gnn.Reserve('getnextnum_test') is: 1000)
		Assert(gnn.Reserve('getnextnum_test') is: 1001)
		gnn.PutBack(1000, 'getnextnum_test')
		Assert(gnn.Reserve('getnextnum_test') is: 1000)
		Assert(gnn.Reserve('getnextnum_test') is: 1002)

		gnn.Date = gnn.Date.Plus(seconds: 200)

		gnn.Confirm(1002, 'getnextnum_test')
		gnn.Renew(1000, 'getnextnum_test')

		gnn.Date = gnn.Date.Plus(seconds: 200)

		// 1001 should have expired so it will be re-used
		Assert(gnn.Reserve('getnextnum_test') is: 1001)
		Assert(gnn.Reserve('getnextnum_test') is: 1003)

		gnn.Date = gnn.Date.Plus(seconds: 200)

		// now 1000 should have expired so it will be re-used
		Assert(gnn.Reserve('getnextnum_test') is: 1000)
		Assert(gnn.Reserve('getnextnum_test') is: 1004)
		gnn.Confirm(1000, 'getnextnum_test')
		gnn.Confirm(1001, 'getnextnum_test')
		gnn.Confirm(1003, 'getnextnum_test')
		gnn.Confirm(1004, 'getnextnum_test')

		Assert(QueryCount('getnextnum_test') is: 1)
		}
	Teardown()
		{
		Suneido.Delete('forceRetryTransaction')
		if TableExists?('getnextnum_test')
			Database('destroy getnextnum_test')
		super.Teardown()
		}
	}