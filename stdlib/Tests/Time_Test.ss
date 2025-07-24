// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(not Time?(''))
		Assert(not Time?('xx'))
		Assert(not Time?('9:30'))
		Assert(Time?('000'))
		Assert(Time?('830'))
		Assert(Time?('0800'))
		Assert(Time?('1200'))
		Assert(Time?('1901'))
		Assert(Time?('2359'))
		Assert(Time?('12'))
		Assert(not Time?('2500'))
		Assert(not Time?('2377'))
		Assert(not Time?('237711'))
		Assert(Time?('0900'))
		Assert(not Time?('12345'))
		Assert(not Time?('9.44'))
		Assert(not Time?('1033.15'))

		// large numbers that can't convert to integer
		Assert(not Time?('4e19'))
		Assert(not Time?('45698732176'))
		Assert(not Time?('-4e19'))
		Assert(not Time?('-45698732176'))
		}
	}