// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rect = #(top: 5, bottom: 10, left: 1, right: 5)
		Assert(POINTinRECT(rect, #(x: 1, y: 5)))
		Assert(not POINTinRECT(rect, #(x: 0, y: 5)))
		Assert(POINTinRECT(rect, #(x: 4, y: 10)))
		Assert(not POINTinRECT(rect, #(x: 10, y: 25)))
		}
	}
