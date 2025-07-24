// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(StateProvName('CA') is: 'California')
		Assert(StateProvName('SK') is: 'Saskatchewan')
		Assert(StateProvName('XX') is: 'XX')
		}
	}