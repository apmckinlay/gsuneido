// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// SuJsWebTest
Test
	{
	Test_one()
		{
		.test(0, #First)
		.test(3, #Last)
		.test(false, #Empty?)
		.test(true, #NotEmpty?)
		.test(true, #Has?, [2])
		.test(false, #Has?, [22])
		.test(4, #Count)
		.test(false, #HasNamed?)

		_next = 0
		seq = Sequence(new .counter)
		seq.HasNamed?()
		seq.Size(named:)
		Assert(_next is: 0)

		.test2(Seq(4))
		.test2(Sequence(new .counter))

		.testboth()
			{|seq|
			Object().Append(seq)
			Assert(not seq.Instantiated?())
			}
		.testboth()
			{|seq|
			Assert(not seq.Instantiated?())
			}
		}
	test(expected, method, args = #())
		{
		seq = Seq(4)
		Assert(seq[method](@args) is: expected)
		Assert(not seq.Instantiated?())
		seq = Sequence(new .counter)
		Assert(seq[method](@args) is: expected)
		Assert(not seq.Instantiated?())
		}
	test2(seq)
		{
		ob = Object()
		for i in seq
			for j in seq
				ob.Add([i, j])
		Assert(ob isSize: 16)
		Assert(ob is: #([0,0], [0,1], [0,2], [0,3], [1,0], [1,1], [1,2], [1,3],
			[2,0], [2,1], [2,2], [2,3], [3,0], [3,1], [3,2], [3,3]))
		}
	testboth(block)
		{
		block(Seq(4))
		block(Sequence(new .counter))
		}
	counter: class
		{
		i: 0
		Next()
			{
			try ++_next
			return .i < 4 ? .i++ : this
			}
		Dup()
			{ return new (.Base()) }
		Infinite?()
			{ return false }
		}
	Test_Count()
		{
		seq = Seq(4)
		Assert(seq.Count() is: 4)
		Assert(not seq.Instantiated?())

		seq = Seq(4)
		Assert(seq.Count(9) is: 0)
		Assert(not seq.Instantiated?())

		seq = Seq(4)
		Assert(seq.Count(2) is: 1)
		Assert(not seq.Instantiated?())
		}
	Test_Without()
		{
		Assert(Fibonaccis().Take(10)
			is: #(0, 1, 1, 2, 3, 5, 8, 13, 21, 34))
		Assert(Fibonaccis().Without(5, 3).Take(10)
			is: #(0, 1, 1, 2, 8, 13, 21, 34, 55, 89))
		}
	}