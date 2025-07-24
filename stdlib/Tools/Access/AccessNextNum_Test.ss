// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.dataTable = .MakeTable('(num, test_TS) key(num)')
		for i in ..20
			QueryOutput(.dataTable, [num: i])
		}

	Test_skip_existing_nextnum()
		{
		table = .MakeNextNum(num: 1000)
		ob = FakeObject(
			AccessNextNum_isDuplicate: function (nextnum)
				{ return nextnum < 1005 },
			AccessNextNum_logNextNumWarning: function (msg)
				{ return msg }
			)
		ob.AccessNextNum_nextNum = [:table, table_field: 'num', field: 'num']
		ob.AccessNextNum_base_query = ""
		ob.AccessNextNum_nextNumClass = GetNextNum
		ob.AccessNextNum_nextNumAttempts = 100

		Assert(ob.Eval(AccessNextNum.AccessNextNum_nextnum, Record()) is: 1005)
		}

	Test_main()
		{
		.nextNumTable = .MakeNextNum(num: 1)
		parent = .parentMock()
		watch = .WatchTable('suneidolog')
		cl = AccessNextNum
			{ AccessNextNum_logNextNumWarning(msg /*unused*/) { } }

		nextNum = [table: .nextNumTable, table_field: 'num', field: 'num']
		accessNextNum1 = cl(parent, nextNum)
		accessNextNum2 = cl(parent, nextNum)
		accessNextNum3 = cl(parent, nextNum)

		.testStart = Timestamp()
		.assertNext(1)

		// Reserve a number, should get 20 (the starting number)
		// 2 is now the next available
		accessNextNum1.SetData(rec1 = [], newrec:)
		Assert(rec1.num is: 20)
		.assertReserved(20)
		.assertNext(21)
		Assert(.GetWatchTable(watch) isSize: 0, msg: 'Test 1')

		// Reserve another number, should get 21 (the second number)
		// 3 is now the next available
		accessNextNum2.SetData(rec2 = [], newrec:)
		Assert(rec2.num is: 21)
		.assertReserved(21)
		.assertNext(22)
		Assert(.GetWatchTable(watch) isSize: 0, msg: 'Test 2')

		// Putback rec1's number (20), new instance attempts a reservation,
		// 20 is reserved, 23 is still the next available
		accessNextNum1.PutBack()
		.assertExpired(20)
		accessNextNum3.SetData(rec3 = [], newrec:)
		Assert(rec3.num is: 20)
		.assertReserved(20)
		.assertNext(22)
		Assert(.GetWatchTable(watch) isSize: 0, msg: 'Test 3')

		// Instance attempts to renew its reservation
		parent.When.GetRecordControl().Return(.controlMock(rec2.num))
		nextNumBefore = Query1(.nextNumTable, num: rec2.num)
		Assert(accessNextNum2.Renew())
		Assert(Query1(.nextNumTable, num: rec2.num).getnextnum_reserved_till
			greaterThan: nextNumBefore.getnextnum_reserved_till)
		Assert(.GetWatchTable(watch) isSize: 0, msg: 'Test 4')

		// Instance confirms reservation is used, removing it from the next num table
		parent.When.GetData().Return(rec2)
		accessNextNum2.Confirm()
		Assert(Query1(.nextNumTable, num: rec2.num) is: false)
		Assert(.GetWatchTable(watch) isSize: 0, msg: 'Test 5')
		}

	parentMock()
		{
		parent = Mock()
		parent.When.AlertWarn([anyArgs:]).Do({ })
		parent.When.GetBaseQuery().Return(.dataTable)
		parent.When.Delay([anyArgs:]).Return(false)
		parent.When.EditMode?().Return(true)
		parent.When.NewRecord?().Return(true)
		return parent
		}

	controlMock(num)
		{
		childControl = Mock()
		childControl.When.Get().Return(num)

		control = Mock()
		control.When.FindControl([anyArgs:]).Return(childControl)
		return control
		}

	assertNext(num)
		{
		Assert(nextNum = Query1(.nextNumTable, :num) isnt: false)
		Assert(nextNum.getnextnum_reserved_till is: Date.Begin())
		}

	assertExpired(num)
		{
		Assert(nextNum = Query1(.nextNumTable, :num) isnt: false)
		Assert(nextNum.getnextnum_reserved_till lessThan: Timestamp())
		}

	assertReserved(num)
		{
		Assert(nextNum = Query1(.nextNumTable, :num) isnt: false)
		Assert(nextNum.getnextnum_reserved_till greaterThan: .testStart)
		}

	Test_Renew_next_num_is_always_used()
		{
		.nextNumTable = .MakeNextNum(num: 1)
		parent = .parentMock()
		cl = AccessNextNum
			{ AccessNextNum_logNextNumWarning(msg) { SuneidoLog(msg) } }
		sulogWatch = .WatchTable('suneidolog')

		nextNum = [table: .nextNumTable, table_field: 'num', field: 'num']
		accessNextNum = cl(parent, nextNum)
		accessNextNum.AccessNextNum_nextNumAttempts = 5
		accessNextNum.AccessNextNum_renewAttempts = 3

		// Attempts to reserve a number, finds that the next several numbers are used.
		// Hits limit before finding an unused number, no number is reserved.
		.testStart = Timestamp()
		accessNextNum.SetData(rec = [], newrec:)
		Assert(rec.num is: 5)
		Assert(Query1(.nextNumTable, num: 5) is: false)
		.assertNext(6)
		Assert(.GetWatchTable(sulogWatch).Last().sulog_message
			is: 'AccessControl.nextnum - had to skip 5 numbers to get: 5')
		parent.Verify.Never().GetData()
		parent.Verify.Never().EditMode?()

		// Attempts to Renew and find a number which is unused.
		// Fails to do find a new number and will not attempt again
		parent.When.GetData().Return(rec)
		Assert(accessNextNum.Renew() matches: 'Please fix your next .* before saving')
		Assert(rec.num is: 15)
		parent.Verify.Times(2).GetData()
		parent.Verify.EditMode?()

		// Ensure it hits the earlier return and does not attempt again
		Assert(accessNextNum.Renew())
		Assert(rec.num is: 15)
		parent.Verify.EditMode?()

		// Ensure we do not run further than the early return as we no longer have a
		// reserved number
		accessNextNum.Confirm()
		parent.Verify.Times(2).GetData()
		}

	Test_Renew_finds_unused_next_num()
		{
		.nextNumTable = .MakeNextNum(num: 1)
		parent = .parentMock()
		cl = AccessNextNum
			{ AccessNextNum_logNextNumWarning(msg) { SuneidoLog(msg) } }
		sulogWatch = .WatchTable('suneidolog')
		nextNum = [table: .nextNumTable, table_field: 'num', field: 'num']
		accessNextNum = cl(parent, nextNum)
		accessNextNum.AccessNextNum_nextNumAttempts = 10
		accessNextNum.AccessNextNum_renewAttempts = 3

		// Attempts to reserve a number, finds that the next several numbers are used.
		// Hits limit before finding an unused number, no number is reserved.
		.testStart = Timestamp()
		accessNextNum.SetData(rec = [], newrec:)
		Assert(rec.num is: 10)
		Assert(Query1(.nextNumTable, num: 10) is: false)
		.assertNext(11)
		Assert(.GetWatchTable(sulogWatch).Last().sulog_message
			is: 'AccessControl.nextnum - had to skip 10 numbers to get: 10')
		parent.Verify.Never().EditMode?()
		parent.Verify.Delay([anyArgs:])
		parent.Verify.Never().AlertWarn([anyArgs:])

		// Attempts to Renew and find a number which is unused.
		// Finds a unused number, reserves it and continues as normal
		parent.When.GetRecordControl().Return(.controlMock(20))
		parent.When.GetData().Return(rec)
		Assert(accessNextNum.Renew())
		Assert(rec.num is: 20)
		.assertNext(21)
		logs = .GetWatchTable(sulogWatch)
		Assert(logs[logs.Size() - 2].sulog_message // second last log
			is: 'renew_next_num while looped 2 times')
		parent.Verify.AlertWarn('Next Number',
			'Another user has taken 10. You have been assigned 20')
		Assert(.GetWatchTable(sulogWatch).Last().sulog_message
			is: 'WARNING: Next Number started at: 10, skipped to: 20')
		parent.Verify.EditMode?()
		parent.Verify.GetData()
		parent.Verify.Times(2).Delay([anyArgs:])

		// Ensure it renews and maintains its reservation
		Assert(accessNextNum.Renew())
		Assert(rec.num is: 20)
		parent.Verify.Times(2).EditMode?()
		parent.Verify.Times(3).Delay([anyArgs:])
		// Ensuring it doesn't reach this point again (AccessNextNum.attemptRenew)
		// Hitting this again would mean its looking for a unused number
		// (which it already has)
		parent.Verify.GetData()
		parent.Verify.AlertWarn([anyArgs:])

		// Ensure we remove our reservation, and that the next number is correctly set
		Assert(Query1(.nextNumTable, num: rec.num) isnt: false)
		accessNextNum.Confirm()
		.assertNext(21)
		Assert(Query1(.nextNumTable, num: rec.num) is: false)
		parent.Verify.Times(2).GetData()
		parent.Verify.AlertWarn([anyArgs:])
		}
	}
