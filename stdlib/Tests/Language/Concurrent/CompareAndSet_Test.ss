// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		ob = Object()
		Assert(ob.CompareAndSet(#x, 123))
		Assert(ob.x is: 123)
		Assert(not ob.CompareAndSet(#x, 123))
		Assert(ob.CompareAndSet(#x, 456, 123))
		Assert(ob.x is: 456)
		Assert(not ob.CompareAndSet(#x, 789, 123))
		Assert(ob.x is 456)
		}
	}