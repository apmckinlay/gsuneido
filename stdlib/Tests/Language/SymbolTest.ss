// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_basic()
		{
		sym = #symbol
		Assert(sym is: "symbol")
		ob = #(symbol: 123)
		Assert(ob[sym] is: ob.symbol)
		}
	Test_symbol_call()
		{
		c = class
			{
			F(@args)
				{ return Object(this, args) }
			}
		ob = c()
		Assert(ob.F(1, 2, a: 3, b: 4) is: Object(ob, #(1, 2, a: 3, b: 4)))
		Assert(#F(ob, 1, 2, a: 3, b: 4) is: Object(ob, #(1, 2, a: 3, b: 4)))

		Assert(#Size("hello") is: 5) // equivalent to "hello".Size()
		Assert(#Size(#(1,2,3)) is: 3)
		}
	}