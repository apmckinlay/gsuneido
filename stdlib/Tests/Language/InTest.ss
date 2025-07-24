// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(2 in () is: false)
		Assert(2 in (2))
		Assert(1 in (1,2,3))
		Assert(2 in (1,2,3))
		Assert(3 in (1,2,3))
		Assert(22 in (1,2,3) is: false)
		x = 0
		Assert(-1 in (x-1, x, x+1))
		Assert(0 in (x-1, x, x+1))
		Assert(+1 in (x-1, x, x+1))
		}
	Test_not()
		{
		Assert(1 in (1,2,3))
		Assert(1 not in (1,2,3) is: false)
		Assert(1 in (2,3) is: false)
		Assert(1 not in (2,3))
		}
	}