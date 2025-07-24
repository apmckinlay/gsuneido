// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Count()
		{
		Assert(Nested.Count(#()) is: 0)
		Assert(Nested.Count(#(123)) is: 1)
		Assert(Nested.Count(#(1, 2, 3)) is: 3)
		Assert(Nested.Count(#(1, (2, (3, 4)), 5)) is: 7)
		}
	Test_Visit()
		{
		test = {|tree, expected|
			actual = Object()
			Nested.Visit(tree, {|x| actual.Add(x) })
			Assert(actual is: expected)
			}
		test(#(1), #(1))
		test(#(1, 2, 3), #(1, 2, 3))
		test(#(1, (2, (3, 4)), 5), #(1, (2, (3, 4)), 2, (3, 4), 3, 4, 5))
		}
	}