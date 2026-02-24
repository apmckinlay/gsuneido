// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		f = Bind(Object)
		Assert(f() is: #())
		Assert(f(1, 2, a: 3, b: 4) is: #(1, 2, a: 3, b: 4))

		f = Bind(Object, 1, 2, a: 3, b: 4)
		Assert(f() is: #(1, 2, a: 3, b: 4))
		Assert(f(33, 44, b: 5) is: #(1, 2, 33, 44, a: 3, b: 5))
		}
	}