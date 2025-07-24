// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		addon = Addon(Date, true)
		Assert(addon.Parent is: Date)
		Assert(addon.Begin() is: Date.Begin())
		}
	Test_options()
		{
		addon = Addon(Date, #(one: 1, Two: 2))
		Assert(addon.One is: 1)
		Assert(addon.Two is: 2)
		}
	}