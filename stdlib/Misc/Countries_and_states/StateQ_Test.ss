// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(State?('') is: false)
		Assert(State?('~') is: false)
		for p in StateCodes
			Assert(State?(p) is: true)
		for p in ProvinceCodes
			Assert(State?(p) is: false)
		}
	}