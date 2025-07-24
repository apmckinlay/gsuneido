// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		ob = QuarterToDateRange('1', '2018')
		Assert(ob.start is: #20180101)
		Assert(ob.end is: #20180331)

		ob = QuarterToDateRange(3, '2005')
		Assert(ob.start is: #20050701)
		Assert(ob.end is: #20050930)

		ob = QuarterToDateRange(4)
		Assert(ob.start is: Date(Date().Year() $ '1001'))
		Assert(ob.end is: Date(Date().Year() $ '1231'))
		}
	}