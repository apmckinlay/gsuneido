// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		targetStr = 'function()
	{
	test1 + 3 + test2(test3[test4].test5)
	if test6
		test7
	}'
		target = Tdop(targetStr)

		// non-existent parttern
		parttern = Tdop('nonexistent', type: 'expression')
		res = TdopSearch(target, parttern)
		Assert(res is: false)

		// single number
		parttern = Tdop('3', type: 'expression')
		res = TdopSearch(target, parttern)
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(targetStr, res[0], '3')

		parttern = Tdop('3', type: 'expression')
		res = TdopSearch(target, parttern, 30)
		Assert(res is: false)

		// complex parttern
		parttern = Tdop('test2(test3[test4].test5)', type: 'expression')
		res = TdopSearch(target, parttern)
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(targetStr, res[0], 'test2(test3[test4].test5)')

		// single placeholder
		parttern = Tdop('a', type: 'expression')
		res = TdopSearch(target, parttern)
		Assert(res isnt: false)
		Assert(res isSize: 2)
		.compare(targetStr, res[0], targetStr)
		.compare(targetStr, res[1], targetStr)

		// complex expression with multiple placeholders
		parttern = Tdop('a+test2(b)', type: 'expression')
		res = TdopSearch(target, parttern)
		Assert(res isnt: false)
		Assert(res isSize: 3)
		.compare(targetStr, res[0], 'test1 + 3 + test2(test3[test4].test5)')
		.compare(targetStr, res[1], 'test1 + 3')
		.compare(targetStr, res[2], 'test3[test4].test5')

		// statement parttern with multiple placeholders
		parttern = Tdop('if a b', type: 'statement')
		res = TdopSearch(target, parttern)
		Assert(res isnt: false)
		Assert(res isSize: 3)
		.compare(targetStr, res[0], 'if test6
		test7')
		.compare(targetStr, res[1], 'test6')
		.compare(targetStr, res[2], 'test7')
		}

	Test_search_direction()
		{
		targetStr = 'function()
	{
	test1 + test2
	test3 * test4
	test5 /*test*/ + test6 + test7
	}'
		target = Tdop(targetStr)
		parttern = Tdop('a + b', type: 'expression')

		res = TdopSearch(target, parttern)
		.compare(targetStr, res[0], 'test1 + test2')

		pos = res[0][0]
		length = res[0][1]
		res = TdopSearch(target, parttern, :pos)
		.compare(targetStr, res[0], 'test1 + test2')

		res = TdopSearch(target, parttern, pos: pos + 1)
		.compare(targetStr, res[0], 'test5 /*test*/ + test6 + test7')

		res = TdopSearch(target, parttern, pos: pos + length)
		.compare(targetStr, res[0], 'test5 /*test*/ + test6 + test7')

		pos = res[0][0]
		length = res[0][1]
		res = TdopSearch(target, parttern, pos: pos + length)
		Assert(res is: false)

		res = TdopSearch(target, parttern, prev:)
		Assert(res is: false)

		res = TdopSearch(target, parttern, pos: targetStr.Size(), prev:)
		.compare(targetStr, res[0], 'test5 /*test*/ + test6 + test7')

		pos = res[0][0]
		length = res[0][1]
		res = TdopSearch(target, parttern, pos: pos + length, prev:)
		.compare(targetStr, res[0], 'test5 /*test*/ + test6 + test7')

		res = TdopSearch(target, parttern, pos: pos + length - 1, prev:)
		.compare(targetStr, res[0], 'test5 /*test*/ + test6')

		res = TdopSearch(target, parttern, :pos, prev:)
		.compare(targetStr, res[0], 'test1 + test2')

		pos = res[0][0]
		length = res[0][1]
		res = TdopSearch(target, parttern, :pos, prev:)
		Assert(res is: false)
		}

	Test_arg_wildcard()
		{
		targetStr = 'function(arg2)
	{
	test("value", arg: 1, :arg2)
	}'
		target = Tdop(targetStr)
		parttern = Tdop('a(b, c, d)', type: 'expression')
		res = TdopSearch(target, parttern)
		.compare(targetStr, res[0], 'test("value", arg: 1, :arg2)')
		.compare(targetStr, res[1], 'test')
		.compare(targetStr, res[2], '"value"')
		.compare(targetStr, res[3], 'arg: 1')
		.compare(targetStr, res[4], ':arg2')
		}

	compare(orginal, location, expected)
		{
		Assert(orginal[location[0]::location[1]] is: expected)
		}
	}