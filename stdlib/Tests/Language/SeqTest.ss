// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_built()
		{
		Assert(Seq(4) is: #(0,1,2,3))
		Assert(#(0,1,2,3) is: Seq(4))
		Assert(Seq(5, 10) is: #(5,6,7,8,9))
		Assert(Seq(2, 8, 2) is: #(2, 4, 6))
		}
	Test_iterate()
		{
		ob = Object()
		for i in .. 4
			ob.Add(i)
		Assert(ob is: #(0,1,2,3))
		}
	Test_join()
		{
		Assert(Seq(4).Join(",") is: "0,1,2,3")
		}
	}