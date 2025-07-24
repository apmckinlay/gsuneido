// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_main()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		addTs = .CommitAdd(svc, svcTable, 'TestRec', 'text', 'TestRec add', 'one')
		modifyTs = .CommitTextChange(svc, svcTable, 'TestRec', 'text changed',
			'minor refactor', 'two')
		modify2Ts = .CommitTextChange(svc, svcTable, 'TestRec', 'text changed again',
			'Issue 1 - modify', 'one')
		modify3Ts = .CommitTextChange(svc, svcTable, 'TestRec', 'text changed 3',
			'Issue 2 - modify', 'one, two')
		stats = SvcStatistics(svc, lib, 'TestRec')

		.testGetData(stats, addTs, modifyTs, modify2Ts, modify3Ts)
		.testGetBasePoint()
		.testGetWeightByDate(stats, modify3Ts)

		Assert(stats.GetContribPercentages()
			is: [[user: 'one', percentage: 75], [user: 'two', percentage: 50]])
		Assert(stats.WeighContributions() is: [['one', 28], ['two', 9.75]])
		}

	testGetData(stats, addTs, modifyTs, modify2Ts, modify3Ts)
		{
		Assert(stats.SvcStatistics_data is: Object(
			[id: 'one', lib_committed: addTs, comment: 'TestRec add'],
			[id: 'two', lib_committed: modifyTs, comment: 'minor refactor'],
			[id: 'one', lib_committed: modify2Ts, comment: 'Issue 1 - modify'],
			[id: 'one, two', lib_committed: modify3Ts, comment: 'Issue 2 - modify']))
		Assert(stats.GetContributors() is: Object('one', 'two'))

		stats.SetName('NonExistent')
		Assert(stats.SvcStatistics_data is: Object())
		Assert(stats.GetContributors() is: Object())
		stats.SetName('TestRec')
		}

	testGetBasePoint()
		{
		fn = SvcStatistics.SvcStatistics_getBasePoint
		Assert(fn('issue 1 - modify', 'one') is: 6)
		Assert(fn('issue 2 - modify', 'one, two') is: 5)
		Assert(fn('modify', 'one') is: 4)
		Assert(fn('modify', 'one, two') is: 3)
		Assert(fn('minor refactor', 'one') is: 2)
		Assert(fn('minor refactor', 'one, two') is: 1)
		}

	testGetWeightByDate(stats, modify3Ts)
		{
		m = stats.SvcStatistics_getWeightByDate
		Assert(m(modify3Ts) is: 1)
		Assert(m(modify3Ts.Minus(days: 20)) is: 1)
		Assert(m(modify3Ts.Minus(days: 40)) is: 0.9)
		Assert(m(modify3Ts.Minus(days: 70)) is: 0.8)
		Assert(m(modify3Ts.Minus(days: 100)) is: 0.7)
		Assert(m(modify3Ts.Minus(days: 130)) is: 0.6)
		Assert(m(modify3Ts.Minus(days: 160)) is: 0.5)
		Assert(m(modify3Ts.Minus(days: 190)) is: 0.4)
		Assert(m(modify3Ts.Minus(days: 220)) is: 0.4)
		Assert(m(modify3Ts.Minus(days: 260)) is: 0.4)
		Assert(m(modify3Ts.Minus(days: 290)) is: 0.4)
		Assert(m(modify3Ts.Minus(days: 320)) is: 0.4)
		Assert(m(modify3Ts.Minus(days: 5000)) is: 0.4)
		}
	}
