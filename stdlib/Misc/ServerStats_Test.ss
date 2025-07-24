// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_uptime()
		{
		m = ServerStats.ServerStats_uptime

		start = asof = #20181204.1500
		Assert(m(start, asof) is: '0 hour(s)')

		start = asof.Plus(minutes: -125)
		Assert(m(start, asof) is: '2 hour(s)')

		start = asof.Plus(days: -5)
		Assert(m(start, asof) is: '5 day(s)')
		}

	Teardown()
		{
		super.Teardown()
		}
	}
