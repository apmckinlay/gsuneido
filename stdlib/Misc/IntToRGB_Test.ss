// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(IntToRGB(0) is: #(r:0, g:0, b:0))
		Assert(IntToRGB(1) is: #(r:1, g:0, b:0))
		Assert(IntToRGB(16744703) is: #(b: 255, g: 128, r: 255)) // teal
		Assert(IntToRGB(16777215) is: #(b: 255, g: 255, r: 255))
		}
	}