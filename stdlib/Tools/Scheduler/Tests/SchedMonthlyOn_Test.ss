// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(SchedMonthlyOn("on StartMonth at 0:00").SchedMonthlyOn_on is: 1)
		Assert(SchedMonthlyOn("on MidMonth at 01:00").SchedMonthlyOn_on is: 15)
		Assert(SchedMonthlyOn("on EndMonth at 01:00").SchedMonthlyOn_on > 27)
		Assert(SchedMonthlyOn("on 7 at 01:00").SchedMonthlyOn_on is: 7)
		Assert(SchedMonthlyOn("on 1 at 11:00") is: false)
		Assert(SchedMonthlyOn("on 40 at 11:00") is: false)
		Assert(SchedMonthlyOn("on 5 at 25:00") is: false)
		Assert(SchedMonthlyOn("on unhandled at 01:00") is: false)
		Assert(SchedMonthlyOn("garbage") is: false)
		}
	}