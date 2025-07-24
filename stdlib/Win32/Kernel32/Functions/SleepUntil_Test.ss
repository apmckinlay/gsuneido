// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		before = Date()
		SleepUntil(before.Plus(milliseconds: 100))
		Assert(Date().MinusSeconds(before) greaterThanOrEqualTo: .1)
		}

	Test_over_max_sleep()
		{
		over24HrsFromNow = Date().Plus(days: 2)
		Assert({ SleepUntil(over24HrsFromNow) } throws: 'sleep time exceeded 24 hours')
		}
	}