// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	gc: GetContributions
		{
		GetContributions_contribs(unused)
			{
			return _contribs
			}
		}
	Test_main()
		{
		.test(#(), #())
		.test(#(#(), #()), #())
		.test(#((a, b), (c, d)), #(a, b, c, d))
		.test(#((a: (1)), (a: (2))), #(a: (1, 2)))
		.test(#((a: (1)), (2)), #(a: (1), 2))
		.test(#((a: (1)), (a: (1))), #(a: (1, 1))) // currently keeps duplicates

		contrib1 = #(a: (x: (9), y: (8)), b: (x: (7), y: (6)))
		contrib2 = #(b: (w: (1), x: (2)), c: (y: (j: 3), z: (k: 4)))
		result = #(a: (x: (9), y: (8)),
			b: (w: (1), x: (7, 2), y: (6)),
			c: (y: (j: 3), z: (k: 4)))
		.test(Object(contrib1, contrib2), result)

		contrib1 = #(a: (b: (c: (d: (1)))), b: (b: (2)))
		contrib2 = #(a: (b: (c: (d: (2)))), b: (b: (3)))
		result = #(a: (b: (c: (d: (1, 2)))), b: (b: (2, 3)))
		.test(Object(contrib1, contrib2), result)

		.testThrows(#((a: 123), (a: 456)))
		.testThrows(#((a: (1)), (a: 2)))
		.testThrows(#((a: 1), (a: (2))))

		contrib1 = #(a: 1)
		contrib2 = #(a: (b: 1, c: 2))
		.testThrows(Object(contrib1, contrib2))
		.testThrows(Object(contrib2, contrib1))

		contrib1 = #(a: (b: 1))
		contrib2 = #(a: (b: (1)))
		.testThrows(Object(contrib1, contrib2))
		.testThrows(Object(contrib2, contrib1))
		}
	test(contribs, result)
		{
		_contribs = contribs
		Assert((.gc).Func("TestContribName") is: result)
		}
	testThrows(contribs)
		{
		_contribs = contribs
		Assert({ (.gc).Func("TestContribName") } throws: "duplicate")
		}

	Test_functions()
		{
		f0 = function () { return #()}
		f1 = function () { return #(a, b) }
		f2 = function () { return #(c, d) }
		.test(Object(f0), #())
		.test(Object(f0, #()), #())
		.test(Object(f0, f0), #())
		.test(Object(f1, f2), #(a, b, c, d))
		.test(Object(#(a, b), f2), #(a, b, c, d))

		f3 = function () { return #(a: (1)) }
		f4 = function () { return #(a: (2)) }
		.test(Object(f3, #(a: (2))), #(a: (1, 2)))
		.test(Object(f3, f4), #(a: (1, 2)))

		f5 = function () { return #(2) }
		.test(Object(f3, f5), #(a: (1), 2))
		.test(Object(f3, f3), #(a: (1, 1)))
		}
	}