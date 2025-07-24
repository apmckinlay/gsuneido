// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Drop(Seq(10), 0) is: #(0,1,2,3,4,5,6,7,8,9))
		Assert(Drop(Seq(10), 5) is: #(5,6,7,8,9))
		Assert(Drop(Seq(10), 10) is: #())
		Assert(Drop(Seq(10), 15) is: #())
		Assert(Drop(#(), 1) is: #())

		seq = Seq(10)
		seq2 = Drop(seq, 5)
		Assert(not seq.Instantiated?())
		Assert(not seq2.Instantiated?())
		}
	}