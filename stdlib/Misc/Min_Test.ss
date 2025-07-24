// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert({ Min() } throws: "Min")
		Assert({ Min(@#()) } throws: "Min")
		Assert({ Min(@+1#(1)) } throws: "Min")

		Assert(Min(123) is: 123)
		Assert(Min(@#(123)) is: 123)
		Assert(Min(@+1#(111, 123)) is: 123)

		Assert(Min(a: 123) is: 123)
		Assert(Min(@#(a: 123)) is: 123)

		Assert(Min(2, 3, 1) is: 1)
		Assert(Min(2, 3, a: 1, b: 4) is: 1)
		Assert(Min(2, 1, a: 3, b: 4) is: 1)
		Assert(Min(@#(2, 3, a: 1, b: 4)) is: 1)
		Assert(Min(@+1#(0, 2, 3, a: 1, b: 4)) is: 1)
		}
	}