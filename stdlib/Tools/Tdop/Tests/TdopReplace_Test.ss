// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_replace()
		{
		src = 'function ()
			{
			a + 1
			b * c
			afds + 432
			return b
			}'
		tree = Tdop(src)
		manager = AstWriteManager(src, tree)

		pat = Tdop('a-b', type: 'expression')
		expect = "function ()
			{
			a + 1
			b * c
			afds + 432
			return b
			}"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')")

		pat = Tdop('a+b', type: 'expression')
		expect = "function ()
			{
			Print('a + 1', 'a', '1')
			b * c
			Print('afds + 432', 'afds', '432')
			return b
			}"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')")

		expect = "function ()
			{
			Print('a + 1', 'a', '1')
			b * c
			afds + 432
			return b
			}"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')", #(count: 1))

		expect = "Print('afds + 432', 'afds', '432')"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')", #(from: 42, to: 52))

		expect = "fds + 432"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')", #(from: 43, to: 52))

		expect = "afds + 43"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')", #(from: 42, to: 51))

		expect = "{
			Print('a + 1', 'a', '1')
			b * c
			Print('afds + 432', 'afds', '432')
			return b
			}"
		.check(expect, manager, pat, "Print('\\0', '\\1', '\\2')", #(from: 16, to: 71))

		expect = "{
			Print('\\0', '\\1', '\\2')
			b * c
			Print('\\0', '\\1', '\\2')
			return b
			}"
		.check(expect, manager, pat, "\=Print('\\0', '\\1', '\\2')", #(from: 16, to: 71))
		}

	check(expect, manager, pat, replacement, args = #())
		{
		args = Object(manager.GetNewWriter(), pat, replacement).Append(args)
		Assert(TdopReplace(@args) is: expect)
		}
	}