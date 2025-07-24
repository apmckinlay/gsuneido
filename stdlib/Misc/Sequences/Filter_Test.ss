// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		t = {|unused| true }
		f = {|unused| false }
		odd = { it % 2 is 1 }
		Assert(Filter(#(), t) is: #())
		Assert(Filter(#(), f) is: #())
		Assert(Filter(#(), odd) is: #())

		Assert(Filter(Seq(5), t) is: #(0,1,2,3,4))
		Assert(Filter(Seq(5), f) is: #())
		Assert(Filter(Seq(10), odd) is: #(1, 3, 5, 7, 9))
		}
	}