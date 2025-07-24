// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Province?('') is: false)
		Assert(Province?('~') is: false)
		for p in ProvinceCodes
			Assert(Province?(p) is: true)
		for p in StateCodes
			Assert(Province?(p) is: false)
		}
	}