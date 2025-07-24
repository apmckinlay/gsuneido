// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_for_in()
		{
		src = "hello"
		dst = Object()
		for c in src
			dst.Add(c)
		Assert(dst is: #(h, e, l, l, o))
		}
	Test_for_range()
		{
		s = ""
		for ..5
			s $= 'x'
		Assert(s is: "xxxxx")

		s = ""
		for i in 2..8
			s $= i
		Assert(s is "234567")

		s = ""
		for i in ..5
			s $= i
		Assert(s is "01234")
		}
	Test_for_m_v()
		{
		test = function (x, expected)
			{
			results = Object()
			for m, v in x
				results.Add([m, v])
			results.Sort!()
			Assert(results is: expected)
			}
		test(#(), #())
		test(#(a, b, c), #((0, a), (1, b), (2, c)))
		test(#(a: 3, b: 4, c: 5), #((a, 3), (b, 4), (c, 5)))
		test(class{ X: 3, Y: 4, Z: 5 }, #((X, 3)(Y, 4)(Z, 5)))
		x = class{}() // instance
		x.a = 3
		x.b = 4
		x.c = 5
		test(x, #((a, 3), (b, 4), (c, 5)))
		test(Seq(10, 13), #((0, 10), (1, 11), (2, 12)))
		}
	}