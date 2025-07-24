// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(DaylightSavings?(#20190101))
		Assert(DaylightSavings?(#20190301))
		Assert(DaylightSavings?(#20190309))
		Assert(DaylightSavings?(#20190310) is: false)
		Assert(DaylightSavings?(#20190615) is: false)
		Assert(DaylightSavings?(#20190701) is: false)
		Assert(DaylightSavings?(#20191102) is: false)
		Assert(DaylightSavings?(#20191103))
		Assert(DaylightSavings?(#20191201))
		Assert(DaylightSavings?(#20191231))

		Assert(DaylightSavings?(#20200101))
		Assert(DaylightSavings?(#20200301))
		Assert(DaylightSavings?(#20200307))
		Assert(DaylightSavings?(#20200308) is: false)
		Assert(DaylightSavings?(#20200615) is: false)
		Assert(DaylightSavings?(#20200701) is: false)
		Assert(DaylightSavings?(#20201031) is: false)
		Assert(DaylightSavings?(#20201101))
		Assert(DaylightSavings?(#20201201))
		Assert(DaylightSavings?(#20201231))
		}
	}