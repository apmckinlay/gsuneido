// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Map(#(), { it }) is: #())
		Assert(Map(Seq(5), { it }) is: #(0,1,2,3,4))
		Assert(Map(Seq(5), { it + 5 }) is: #(5,6,7,8,9))

		seq = Seq(10)
		seq2 = Map(seq, { it + 1 })
		Assert(not seq.Instantiated?(), msg: 'seq.Instantiated?')
		Assert(not seq2.Instantiated?(), msg: 'seq2.Instantiated?')
		}
	}