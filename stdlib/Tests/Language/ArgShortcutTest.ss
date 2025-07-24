// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		f = function (@args) { return args }
		a = 1
		b = 2
		Assert(f(0, :a, :b) is: #(0, a: 1, b: 2))
		Assert({ "Print(:x)".Eval() } throws: 'uninitialized')  // .Eval is okay here
		Assert({ "function () { f(:a, :a) }".Compile() }
			throws: 'duplicate argument name: a')
		Assert([0, :a, :b] is: #{0, a: 1, b: 2})

		Assert(f(@[:a, :b]) is: #(a: 1, b: 2))
		Assert(f([:a, :b]) is: #((a: 1, b: 2)))

		Assert(f(a: [:b]) is: #(a: (b: 2)))
		}
	}
