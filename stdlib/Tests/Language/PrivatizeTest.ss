// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		x = new PrivatizeTestHelper
		Assert(x members: #(#PrivatizeTestHelper_a, #b, 'c'))
		}
	}