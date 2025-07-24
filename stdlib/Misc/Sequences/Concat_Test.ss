// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Concat(#()) is: #())
		Assert(Concat(#(), #()) is: #())
		Assert(Concat(#(), #(), Seq(3)) is: #(0, 1, 2))
		Assert(Concat(#(), Seq(3) #(), #(3, 4, 5)) is: #(0, 1, 2, 3, 4, 5))

		x = Seq(3)
		y = Seq(3)
		concat = Concat(x, y)
		Assert(not x.Instantiated?())
		Assert(not y.Instantiated?())
		Assert(not concat.Instantiated?())

		x = Seq(3)
		y = Seq() // infinite
		z = Seq(3)
		Assert(Concat(y).Infinite?())
		Assert(not Concat(x, z).Infinite?())
		Assert(Concat(x, y, z).Infinite?())
		Assert(Concat(y, x, z).Infinite?())
		Assert(Concat(x, z, y).Infinite?())
		}
	}