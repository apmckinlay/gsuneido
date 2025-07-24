// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Nof(0, { 123 }) is: #())
		Assert(Nof(4, { 123 }) is: #(123, 123, 123, 123))
		}
	}