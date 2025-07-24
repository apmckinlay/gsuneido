// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(#(a,b,c).Map2({|i,v| i $ "->" $ v })
			is: #("0->a", "1->b", "2->c"))
		}
	}