// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(ReadableDuration(.123) is: "123 ms")
		Assert(ReadableDuration(1.234) is: "1.23 sec")
		Assert(ReadableDuration(180) is: "3 min")
		Assert(ReadableDuration(34 * 60) is: "34 min")
		Assert(ReadableDuration(3 * 3600) is: "3 hr")
		}
	}