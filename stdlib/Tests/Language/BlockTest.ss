// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// SuJsWebTest
Test
	{
	Test_main()
		{
		b = { 123 }
		Assert(b() is: 123)
		x = 123
		b = { x + 1 }
		Assert(b() is: 124)
		b = { y = 456 }
		b()
		Assert(y is: 456)
		b = {|a, b, c| a + b + c }
		Assert(b(1, 2, 3) is: 6)

		f = function () { n = 0; return { ++n } }
		b = f()
		Assert(b() is: 1)
		Assert(b() is: 2)

		f = function (ob) { n = 0; ob.b = { ++n }; return }
		f(ob = Object())
		Assert((ob.b)() is: 1)
		Assert((ob.b)() is: 2)
		}

	Test_break_continue()
		{
		myfor = function (block)
			{
			for i in ..10
				try
					block(i)
				catch (x, "block:")
					if (x is "block:break")
						return
					// else block:continue ... so continue
			}
		ob = Object()
		myfor()
			{|n|
			if (n is 3)
				continue
			if (n > 4)
				break
			ob.Add(n)
			}
		Assert(ob is: #(0, 1, 2, 4))
		}

	Test_return()
		{
		Assert(.test_return() is: 'done')
		}
	test_return()
		{
		run = function (block) { block() }
		run()
			{
			run()
				{
				return 'done'
				}
			}
		return 'should not get here!'
		}

	Test_return_from_tran()
		{
		// SuJsWebTest Excluded
		table = .MakeTable("(a,b,c) key(a)")
		.test_return_from_tran(table)
		Assert(Query1(table) is: [a:1])
		}
	test_return_from_tran(table)
		{
		Transaction(update:)
			{|t|
			t.QueryOutput(table, [a:1])
			return
			}
		}

	Test_return_not_caught()
		{
		.test_caller({ return 123 })
		}
	test_caller(block)
		{
		try
			block()
		catch (e)
			throw "should not get here"
		}

	Test_shortcuts()
		{
		// NOTE: deliberately using Times and blocks here
		n = 0
		10.Times({ ++n })
		Assert(n is: 10)
		10.Times() { ++n }
		Assert(n is: 20)
		10.Times { ++n }
		Assert(n is: 30)
		}
	Test_it()
		{
		list = #(1, 2, 3)
		log = []
		list.Each {|it| log.Add(it) }
		Assert(log is: list)
		log = []
		list.Each { log.Add(it) }
		Assert(log is: list)
		}
	Test_args()
		{
		b = {|@args| x = args }
		b()
		Assert(x is: #())
		b(1, 2, 3)
		Assert(x is: #(1, 2, 3))
		b(1, 2, a: 3, b: 4)
		ob = #(1, 2, a: 3, b: 4)
		Assert(x is: ob)
		b(@ob)
		Assert(x is: ob)
		}
	Test_params()
		{
		b = { }
		try { b(123); Assert(false) }
			catch (unused, 'too many arguments') { }

		b = {|x| log = [x] }
		b(123)
		Assert(log is: [123])
		try { b(); Assert(false) }
		try { b(12, 34); Assert(false) }
			catch (unused, 'too many arguments') { }
		b(x: 456)
		Assert(log is: [456])
		b(123, y: 456)
		Assert(log is: [123])

		b = {|x, y| log = [x, y] }
		try { b(); Assert(false) }
			catch (unused, 'missing argument') { }
		try { b(123); Assert(false) }
			catch (unused, 'missing argument') { }
		try { b(12, 34, 56); Assert(false) }
			catch (unused, 'too many arguments') { }
		b(123, 456)
		Assert(log is: [123, 456])
		b(y: 123, x: 456)
		Assert(log is: [456, 123])
		b(123, y: 456)
		Assert(log is: [123, 456])
		b(12, y: 34, z: 56)
		Assert(log is: [12, 34])
		}
	Test_Add()
		{
		ob = Object()
		.one(ob)
		Assert(.two(ob) is: 246)
		}
	one(ob)
		{
		a = b = 123
		ob.Add({ a + b})
		}
	two(ob)
		{
		return (ob[0])()
		}
	Test_thread()
		{
		// SuJsWebTest Excluded
		a = b = 123
		Thread({ Assert(a + b is: 246) })
		// NOTE: will throw in other thread, test itself won't fail
		}
	Test_block_params_dont_clash_with_dynamic_variables()
		{
		fn = function (_b = 1) { b }
		Assert(fn() is: 1)

		block = { |b| b }
		_b = 5
		Assert(block(2) is: 2)
		Assert(fn() is: 5)
		Assert(block(3) is: 3)
		}
	}
