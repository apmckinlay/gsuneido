// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Take(#(), 1) is: #())
		Assert(Take(Seq(10), 0) is: #())
		Assert(Take(Seq(10), 5) is: #(0,1,2,3,4))
		Assert(Take(Seq(10), 10) is: #(0,1,2,3,4,5,6,7,8,9))
		Assert(Take(Seq(10), 15) is: #(0,1,2,3,4,5,6,7,8,9))

		seq = Seq(10)
		seq2 = Take(seq, 5)
		Assert(not seq.Instantiated?())
		Assert(not seq2.Instantiated?())
		}
	}