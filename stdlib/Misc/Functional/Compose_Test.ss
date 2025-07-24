// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		add = function (x, y) { x + y }
		neg = function (x) { -x }
		Assert(Compose(add)(123, 456) is: 579)
		Assert(Compose(add, neg)(4, 1) is: -5)
		}
	}