// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Lines()
		{
		s = "one\ntwo\nthree"

		Assert(Lines(s) isSize: 3)
		Assert(Lines(s).Find('two') is: 1)
		Assert(Lines(s) is: #(one, two, three))

		result = Object()
		for v in Lines(s)
			result.Add(v)
		Assert(result is: #(one, two, three))

		x = Lines(s).Iter()
		x.Next()
		Assert(x.Remainder() is: "two\nthree")
		}

	Test_Values()
		{
		ob = #(one, two, three)

		Assert(ob.Values() isSize: 3)
		Assert(ob.Values().Find('two') is: 1)
		Assert(ob.Values() is: #(one, two, three))

		result = Object(0)
		for v in ob.Values()
			result.Add(v)
		Assert(result is: #(0, one, two, three))
		}

	Test_Seq()
		{
		Assert(Seq(1, 4) isSize: 3)
		Assert(Seq(1, 4).Find(2) is: 1)
		Assert(Seq(1, 4) is: #(1, 2, 3))

		result = Object(0)
		for v in Seq(1, 4)
			result.Add(v)
		Assert(result is: #(0, 1, 2, 3))
		}
	}