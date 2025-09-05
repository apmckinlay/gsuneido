// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.statsMem = .TempName()
		}

	Test_main()
		{
		_resTime = 0
		cl = RackPerformanceStats
			{
			New(app, .StatsMem)
				{
				super(app)
				}
			App(@args)
				{
				return args
				}
			RackPerformanceStats_getDuration(unused)
				{
				return _resTime
				}
			}

		cl = new cl('app', .statsMem)

		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 1)

		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 2)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 3)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 0.5
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 4)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 5)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 6)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 1
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 7)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 8)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][0] is: 9)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 1.5
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][1] is: 1)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][1] is: 2)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][1] is: 3)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 2
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][1] is: 4)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][1] is: 5)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][1] is: 6)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 100
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][7] is: 1)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][7] is: 2)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][7] is: 3)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 1000
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][10] is: 1)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][10] is: 2)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][10] is: 3)
		Assert(Suneido[cl.StatsMem] isSize: 20)

		_resTime = 284804
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][19] is: 1)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][19] is: 2)
		cl([env: 'env'])
		Assert(Suneido[cl.StatsMem][19] is: 3)
		Assert(Suneido[cl.StatsMem] isSize: 20)
		}

	Teardown()
		{
		ServerSuneido.DeleteMember(.statsMem)
		super.Teardown()
		}
	}