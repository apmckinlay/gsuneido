// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		date = #(wHour: 10, wMonth: 3,
			wMilliseconds: 380, wSecond: 54, wMinute: 42, wDay: 15,
			wDayOfWeek: 0, wYear: 2009)
		Assert(SystemTimeToSuneidoDate(date) is: #20090315.104254380)
		}
	}