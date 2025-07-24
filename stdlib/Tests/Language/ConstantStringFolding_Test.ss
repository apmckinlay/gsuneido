// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_object()
		{
		x = #{ A: "hello" $ " " $
					"" $ "world" }
		Assert(x.A is: "hello world")
		}
	Test_class()
		{
		c = class
			{
			A: "hello" $ " " $
			   "" $ "world"
			}
		Assert(c.A is: "hello world")
		}
	}