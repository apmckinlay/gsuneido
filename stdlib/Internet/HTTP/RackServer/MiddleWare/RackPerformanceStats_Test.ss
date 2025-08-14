// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.statsMem = .TempName()
		}

	Test_main()
		{
		_resTime = 0.5
		cl = RackPerformanceStats
			{
			App(@args)
				{
				return args
				}
			RackPerformanceStats_getResTime(unused)
				{
				return _resTime
				}
			}

		cl = new cl('app')
		cl.StatsMem = .statsMem

		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 0) is: 1)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 0) is: 2)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 0) is: 3)
		Assert(ServerSuneido.Get(cl.StatsMem) isSize: 1)

		_resTime = 1
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 0) is: 4)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 0) is: 5)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 0) is: 6)
		Assert(ServerSuneido.Get(cl.StatsMem) isSize: 1)

		_resTime = 1.5
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 1) is: 1)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 1) is: 2)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 1) is: 3)
		Assert(ServerSuneido.Get(cl.StatsMem) isSize: 2)

		_resTime = 2
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 1) is: 4)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 1) is: 5)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 1) is: 6)
		Assert(ServerSuneido.Get(cl.StatsMem) isSize: 2)

		_resTime = 100
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 7) is: 1)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 7) is: 2)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 7) is: 3)
		Assert(ServerSuneido.Get(cl.StatsMem) isSize: 3)

		_resTime = 1000
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 10) is: 1)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 10) is: 2)
		cl([env: 'env'])
		Assert(ServerSuneido.GetAt(cl.StatsMem, 10) is: 3)
		Assert(ServerSuneido.Get(cl.StatsMem) isSize: 4)
		}

	Teardown()
		{
		ServerSuneido.DeleteMember(.statsMem)
		super.Teardown()
		}
	}