// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		types = #(false, true, -1, 0, 1, "", "x", #19990304, #20000304)
		Assert(types.Copy().Sort!() is: types)
		}
	}
