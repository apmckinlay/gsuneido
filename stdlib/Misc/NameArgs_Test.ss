// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.test(#(), [], #())
		.test(#(a: 1, b: 2), [1, 2], #(a, b))
		.test(#(a: 1, b: 2), [a: 1, b: 2], #(a, b))
		.test(#(a: 1, b: 2, c: 3), [a: 1, b: 2, c: 3], #(a, b))
		}
	Test_defaults()
		{
		.test(#(a: 1, b: 2), [], #(a, b), #(1, 2))
		.test(#(a: 1, b: 2), [1, 2], #(a, b), #(3, 4))
		.test(#(a: 1, b: 2, c: 3), [1, 2], #(a, b, c), #(3))
		.test(#(a: 1, b: 2, c: 3), [1, c: 3], #(a, b, c), #(2, 4))
		}
	Test_named_override_unnamed()
		{
		.test(#(a: 1, b: 2, c: 4), [1, 2, 3, c: 4], #(a, b, c))
		}
	Test_too_many_arguments()
		{
		Assert({NameArgs([], #(a))} throws:)
		}
	Test_missing_argument()
		{
		Assert({NameArgs([1], #())} throws:)
		}

	test(expect, args, names, defs = #())
		{
		Assert(NameArgs(args, names, defs) is: expect)
		}
	}