// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(SINValid?('593 486 780') is: true)
		Assert(SINValid?('193 456 787') is: true)
		Assert(SINValid?('393 456 783') is: false)

		Assert(SINValid?('508154671I') is: false)
		Assert(SINValid?('508 154 671') is: true)
		}
	}
