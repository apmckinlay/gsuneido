// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(AstSearch('AB', [], ['ABC']) is: #())
		Assert(AstSearch('ABC', [], ['ABC']) is: #((pos: 0, end: 3)))
		Assert(AstSearch('true', [], [true, false]) is: #((pos: 0, end: 4)))
		Assert(AstSearch('1_000', [], [1000]) is: #((pos: 0, end: 5)))
		Assert(AstSearch('1_000', [], [100]) is: #())

		target = 'function()
	{
	test1 + 4 + test2(test3[test4].test5)
	a = 3
	if test6
		test7 + 4
	}'

		// non-existent pattern
		res = AstSearch(target, 'nonexistent')
		Assert(res is: #())

		// single number
		res = AstSearch(target, '3')
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(target, res[0], '3')

		// complex pattern
		res = AstSearch(target, 'test2(test3[test4].test5)')
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(target, res[0], 'test1 + 4 + test2(test3[test4].test5)')

		// test nary
		res = AstSearch(target, 'test1 + 4')
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(target, res[0], 'test1 + 4 + test2(test3[test4].test5)')

		res = AstSearch(target, '4+a')
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(target, res[0], 'test1 + 4 + test2(test3[test4].test5)')

		// complex expression with multiple placeholders
		res = AstSearch(target, 'a+test2(b)')
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(target, res[0], 'test1 + 4 + test2(test3[test4].test5)')

		// multiple matches
		res = AstSearch(target, 'a + 4')
		Assert(res isnt: false)
		Assert(res isSize: 2)
		.compare(target, res[0], 'test7 + 4')
		.compare(target, res[1], 'test1 + 4 + test2(test3[test4].test5)')

		// test findFirst?
		res = AstSearch(target, 'a + 4', findFirst?:)
		Assert(res isnt: false)
		Assert(res isSize: 1)
		.compare(target, res[0], 'test7 + 4')

		// test skip?
		s = 'class
	{
	Skip1()
		{
		Abc()
		}
	Skip2: #(a: "Abc")
	Foo: class
		{
		Foo()
			{
			Abc
			}
		Skip3()
			{
			Abc
			}
		}
	}'
		skipFn? = { |node, unused| node.type is 'Member' and node.key isnt 'Foo' }
		res = AstSearch(s, 'Abc', :skipFn?)
		Assert(res isSize: 1)
		Assert(res has: [pos: 97, end: 100])

		target = 'function () { #(a: 1) }'
		res = AstSearch(target, '#(a: 1)')
		Assert(res isSize: 1)
		.compare(target, res[0], '#(a: 1)')

		target = 'Abc {}'
		res = AstSearch(target, 'Abc')
		Assert(res isSize: 1)
		.compare(target, res[0], 'Abc {}')
		}

	Test_arg_wildcard()
		{
		target = 'function(arg2)
	{
	test("value", arg1: 1, :arg2){}
	}'
		// all wildcards
		res = AstSearch(target, 'a(b, c, d, e)')
		Assert(res isSize: 1)
		.compare(target, res[0], 'test("value", arg1: 1, :arg2){}')

		// test any args
		Assert(.match('test()', 'test(@ANY_ARGS)'))
		Assert(.match('test(1, 2, 3)', 'test(@ANY_ARGS)'))
		Assert(.match('test(1, a: 2, b: 3)', 'test(@ANY_ARGS)'))
		Assert(.match('test(@args)', 'test(@ANY_ARGS)'))

		// test skiping extra named args
		Assert(.match('test(1, 2)', 'test(x, y)'))
		Assert(.match('test(1, 2)', 'test(x)') is: false)
		Assert(.match('test(1, 2)', 'test(x, y, z)') is: false)

		Assert(.match('test(a: 1, b: 2)', 'test(a: x, b: y)'))
		Assert(.match('test(a: 1, b: 2)', 'test(a: x, b: y, c: z)') is: false)
		Assert(.match('test(a: 1, b: 2)', 'test(a: x)'))

		Assert(.match('test(1, a: 2)', 'test(x)'))
		Assert(.match('test(1, a: 2)', 'test(a: x)') is: false)
		Assert(.match('test(1, a: 2)', 'test(a: x, b: y)'))
		Assert(.match('test(1, a: 2)', 'test(x, y)'))
		Assert(.match('test(1, a: 2)', 'test(x, a: y)'))
		Assert(.match('test(1, a: 2)', 'test(x, b: y)') is: false)
		Assert(.match('test(1, a: 2)', 'test(x, y, a: z)') is: false)
		Assert(.match('test(1, a: 2)', 'test(x, a: y, b: z)') is: false)
		Assert(.match('test(1, a: 2)', 'test(a: x, b: y, c: z)') is: false)

		// test block
		Assert(.match('test(){}', 'test()'))
		Assert(.match('test(){}', 'test(){}'))
		Assert(.match('test(){}', 'test(a)'))
		Assert(.match('test(){}', 'test(block: a)'))

		// not wildcard
		Assert(.match('test("value")', 'test("value 2")') is: false)
		Assert(.match('test("value")', 'test(arg: "value")'))
		Assert(.match('test("value")', 'test(arg: "value 2")') is: false)
		Assert(.match('test(arg: "value")', 'test(arg: "value")'))
		Assert(.match('test(arg: "value")', 'test(arg: "value 2")') is: false)
		Assert(.match('test(arg: "value")', 'test(arg1: "value")') is: false)
		Assert(.match('test(arg: "value")', 'test("value")'))
		Assert(.match('test(arg: "value")', 'test("value 2")') is: false)

		// wildcard
		Assert(.match('test("value")', 'test(a)'))
		Assert(.match('test("value")', 'test(arg: a)'))
		Assert(.match('test(arg: "value")', 'test(a)'))
		Assert(.match('test(arg: "value")', 'test(arg: a)'))
		Assert(.match('test(arg: "value")', 'test(arg1: a)') is: false)
		Assert(.match('test("value, "arg: "value")', 'test(a, b)'))
		Assert(.match('test("value, "arg: "value")', 'test(a, b, c)') is: false)
		Assert(.match('test("value, "arg: "value")', 'test(a)'))

		// handles when block is in the params
		Assert(.match('test("value"){}', 'test(block: b, arg: "value")'))
		Assert(.match('test("value"){}', 'test(arg: "value", block: b)'))
		Assert(.match('test("value", :arg2, arg1: 1){}',
			'test(a, arg2: b, arg1: 1, block: {})'))
		Assert(.match('test("value", :arg2, arg1: 1){}',
			'test(a, block: {}, arg2: b, arg1: 1)'))
		Assert(.match('test("value", :arg2, arg1: 1){}',
			'test(a, {}, arg2: b, arg1: 1)'))
		}

	Test_mem_wildcard()
		{
		Assert(.match('.aaa = .bbb', '.ccc') is: false)
		Assert(.match('.aaa = .bbb', '.aaa'))
		Assert(.match('.aaa = .bbb', '.c'))
		Assert(.match('.aaa = .bbb', '.c = .d'))

		Assert(.match('test[a + b] = 123', '.c') is: false)
		Assert(.match('test[a + b] = 123', 'test.c'))
		Assert(.match('test[.abc + test["abc"]]', 'a[.b + c.d]'))

		target = 'class
	{
	CallClass(fn = function(a) { .test $ a })
		{
		.Test = this[fn(a)]
		}
	}'
		Assert(AstSearch(target, '.a').Size() is: 3)
		}

	Test_multiple_searches()
		{
		target = 'function ()
			{
			c = Object()
			a = b + 1
			c[a + 1] = 2
			return c[a + 1] + 1
			}'

		Assert(AstSearch(target, 'x + 1') isSize: 4)
		Assert(AstSearch(target, ['x + 1']) isSize: 4)
		Assert(AstSearch(target, 'x + 2') isSize: 0)
		Assert(AstSearch(target, 'x = y') isSize: 3)
		Assert(AstSearch(target, ['x + 1', 'x + 2', 'x = y']) isSize: 7)
		Assert(AstSearch(target, ['x + 1', 'x + 2', 'x = y'], findFirst?:) isSize: 1)

		Assert(AstSearch(target, []) isSize: 0)
		Assert(AstSearch(target, ['x + 2', 'x + 3']) isSize: 0)

		Assert(AstSearch(target, 'x +') hasPrefix: 'Parse Search text "x +" error')
		}

	Test_search_constant()
		{
		target = `function (xyz, b = 'abc', c = #(abc, xyz: 'abc'), d = 'abc',
			test = false)
			{
			a = "abc" $ "xxyz"
			c["abc"] = c.abc $ a $ d[1::1]
			switch (xyz)
				{
			case 'abc':
				return #(#abc, xyz: 'abc')
			case false:
				Xyz(xyz: false)
			default:
				throw 'abc'
				}
			}`

		Assert(AstSearch(target, [], ['abc']) isSize: 11)
		Assert(AstSearch(target, [], [1]) isSize: 2)
		Assert(AstSearch(target, [], [false]) isSize: 3)
		Assert(AstSearch(target, [], ['xyz']) isSize: 0)
		Assert(AstSearch(target, [], ['abc', 1, false, 'xyz']) isSize: 16)

		Assert(AstSearch(target, ['Xyz'], ['abc', 1, false, 'xyz']) isSize: 17)
		Assert(AstSearch(target, ['Xyz', 'xyz'], ['abc', 1, false, 'xyz']) isSize: 18)

		Assert(AstSearch('#(abc, 123, true, #20240112)', [], ['abc']) isSize: 1)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], ['ABC']) isSize: 0)
		Assert(AstSearch('#(abc, 123, true, #20240112)', ['abc'], []) isSize: 0)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], [123]) isSize: 1)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], [12]) isSize: 0)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], [true]) isSize: 1)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], [false]) isSize: 0)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], [#20240112]) isSize: 1)
		Assert(AstSearch('#(abc, 123, true, #20240112)', [], [#20240113]) isSize: 0)

		Assert(AstSearch('#(s: abc)', [], ["abc"]) isSize: 1)
		Assert(AstSearch('#(s: abc)', [], ["s"]) isSize: 0)

		c = 0
		Assert(AstSearch('#(aba, 123, true, #20240112, #20240114 abb: "abc")', [],
			[{ |code| c++; String?(code) and code.Prefix?(#ab) }]) isSize: 2)
		Assert(c is: 6)

		Assert(AstSearch('#(aba, 123, true, #20240112, abb: "abc")', [],
			[{ |code| String?(code) and code.Prefix?(#ab) },
				true,
				{ |code| Date?(code) and code < #20240113 }]) isSize: 4)

		}

	Test_ForIn()
		{
		target = `function ()
			{
			f1 = function (ob) { for i in ob i }
			f2 = function (from, to) { for i in from..to i }
			f3 = function (n) { for ..n i }
			}`

		Assert(AstSearch(target, 'function (ob) { for i in ob i}') isSize: 1)
		Assert(AstSearch(target, 'function (from, to) { for i in from..to i }') isSize: 1)
		Assert(AstSearch(target, 'function (from, to) { for i in from..from i }')
			isSize: 0)
		Assert(AstSearch(target, 'function (n) { for ..n i }') isSize: 1)
		Assert(AstSearch(target, 'function (n) { for ..n+1 i }') isSize: 0)
		}

	match(s, search)
		{
		target = 'function(arg2) { ' $ s $ ' }'
		return AstSearch(target, search).NotEmpty?()
		}

	compare(orginal, location, expected)
		{
		Assert(orginal[location.pos..location.end] is: expected)
		}

	Test_GetHint()
		{
		fn = AstSearch.GetHint

		Assert(fn('a + 1') is: false)
		Assert(fn('a + b') is: false)
		Assert(fn('1 + 1') is: false)
		Assert(fn('.a = false') is: false)

		Assert(fn('test + 1') is: 'test')
		Assert(fn('.test') is: 'test')
		Assert(fn('Query1(a)') is: 'Query1')
		Assert(fn('Query1(@ANY_ARGS)') is: 'Query1')
		Assert(fn('QueryFirst("stdlib sort name")') is: 'stdlib sort name')
		Assert(fn('#(123, abc: a)') is: 'abc')
		Assert(fn('#(123, abc: "abcd")') is: 'abcd')
		Assert(fn('#(("12345"), abc: "abcd")') is: '12345')
		}

	Test_MultiAssign()
		{
		target = `function ()
			{
			v1, v2 = TestFn()
			v3 = TestFn()
			v4, v5 = TestFn(args)
			}`

		Assert(AstSearch(target, 'v1, v2 = TestFn()') isSize: 1)
		Assert(AstSearch(target, 'a, b = c()') isSize: 1)
		Assert(AstSearch(target, 'a, b = c(@ANY_ARGS)') isSize: 2)
		Assert(AstSearch(target, 'a = b') isSize: 1)
		Assert(AstSearch(target, 'TestFn') isSize: 3)

		target = `function ()
			{
			fn0 = function ()
				{
				return TestFn(1), TestFn(2), TestFn(3)
				}
			fn1 = function ()
				{
				return
				}
			fn2 = function ()
				{
				return TestFn(1)
				}
			}`
		Assert(AstSearch(target, 'TestFn(@ANY_ARGS)') isSize: 4)
		Assert(AstSearch(target, 'function () { return }') isSize: 1)
		Assert(AstSearch(target, 'function () { return a }') isSize: 1)
		Assert(AstSearch(target, 'function () { return a, b, c }') isSize: 1)
		}
	}
