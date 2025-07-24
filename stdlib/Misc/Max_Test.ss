// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert({ Max() } throws: "Max")
		Assert({ Max(@#()) } throws: "Max")
		Assert({ Max(@+1#(1)) } throws: "Max")

		Assert(Max(123) is: 123)
		Assert(Max(@#(123)) is: 123)
		Assert(Max(@+1#(999, 123)) is: 123)

		Assert(Max(a: 123) is: 123)
		Assert(Max(@#(a: 123)) is: 123)

		Assert(Max(2, 3, 1) is: 3)
		Assert(Max(2, 3, a: 1) is: 3)
		Assert(Max(@#(2, 3, a: 1)) is: 3)

		Assert(Max(2, 3, a: 4, b: 1) is: 4)
		Assert(Max(@#(2, 3, a: 4, b: 1)) is: 4)
		Assert(Max(@+1#(9, 2, 3, a: 4, b: 1)) is: 4)
		}
	}