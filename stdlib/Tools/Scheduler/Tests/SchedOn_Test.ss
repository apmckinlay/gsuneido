// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_parse()
		{
		Assert(SchedOn("garbage") is: false)
		Assert(SchedOn("on Fri garbage") is: false)

		sched = SchedOn("on Fri at 20:00")
		Assert(sched.SchedOn_on is: 'Fri')
		}
	}