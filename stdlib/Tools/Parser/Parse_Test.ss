// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_alts()
		{
		f = Parse.Parse_alts
		Assert(f(#()) is: #(()))
		Assert(f(#(12 34)) is: #((12 34)))
		Assert(f(#(12 34 or 56)) is: #((12 34), (56)))
		}
	}