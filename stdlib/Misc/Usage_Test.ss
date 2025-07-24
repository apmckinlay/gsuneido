// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.TearDownIfTablesNotExist('usages')
		.type = .TempName()
		.type2 = .TempName()

		today = Date().NoTime()
		Assert(Usage.Get(.type, today) is: 0)
		Assert(Usage.GetDailyTotals(.type, today) is: 0)

		Usage.Increment(.type, today)
		Assert(Usage.Get(.type, today) is: 1)
		Assert(Usage.GetDailyTotals(.type, today) is: 1)

		Usage.Increment(.type, today)
		Assert(Usage.Get(.type, today) is: 2)
		Assert(Usage.GetDailyTotals(.type, today) is: 2)

		yesterday = today.Minus(days: 1)
		Usage.Increment(.type, yesterday)
		Assert(Usage.Get(.type, today) is: 2)
		Assert(Usage.GetDailyTotals(.type, today) is: 2)
		Assert(Usage.Get(.type, yesterday) is: 1)
		Assert(Usage.GetDailyTotals(.type, yesterday) is: 1)

		Usage.Increment(.type, yesterday)
		Assert(Usage.Get(.type, yesterday) is: 2)
		Assert(Usage.GetDailyTotals(.type, yesterday) is: 2)

		Usage.Increment(.type2, yesterday)
		Assert(Usage.Get(.type2, yesterday) is: 1)
		Assert(Usage.GetDailyTotals(.type2, yesterday) is: 1)
		}

	Teardown()
		{
		Usage.Remove(.type)
		Usage.Remove(.type2)
		super.Teardown()
		}
	}