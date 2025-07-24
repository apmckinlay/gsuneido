// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		fn = ExcludeLowerFields
		Assert(fn(#()) is: #())
		Assert(fn(#(a, b, c)) is: #())
		Assert(fn(#(a, b, c), #(a, b)) equalsSet: #(a, b))
		Assert(fn(#(a, b, c, c_lower!), #(a, b)) equalsSet: #(a, b, c_lower!))
		Assert(fn(#(a, b_lower!, c, c_lower!), #(a, b))
			equalsSet: #(a, b, b_lower!, c_lower!))
		}
	}
