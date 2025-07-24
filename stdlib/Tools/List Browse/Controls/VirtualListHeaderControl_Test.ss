// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_fitWindowSize()
		{
		fit = VirtualListHeaderControl.VirtualListHeaderControl_fitWindowSize

		fit(85, widths = Object(50, 50, 30, 30, 10), Mock())
		Assert(widths is: #(25, 25, 15, 15, 5))

		fit(260, widths = Object(100, 100, 60, 40, 30, 50, 50), Mock())
		Assert(widths is: #(60, 60, 36, 24, 18, 30, 30))

		fit(42, widths = Object(20, 22, 40, 30), Mock())
		Assert(widths is: #(7, 8, 15, 11))
		}
	}