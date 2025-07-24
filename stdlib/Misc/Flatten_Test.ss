// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Flatten() is: #())
		Assert(Flatten(#()) is: #())
		Assert(Flatten(#(1, 2)) is: #(1, 2))
		Assert(Flatten(#(1, 2) #(3, 4)) is: #(1, 2, 3, 4))
		Assert(Flatten(#(1, 2), #(3, 4), #()) is: #(1, 2, 3, 4))
		}
	}