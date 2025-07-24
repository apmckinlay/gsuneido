// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_dynamic_implicit()
		{
		_q = 123
		Assert(function(p){ p }(123) is: 123)
		Assert(function(_p){ p }(123) is: 123)
		Assert(function(_p){ p }(p: 123) is: 123)
		Assert({ function(_p){ p }() }, throws: "missing argument")
		Assert(function(_q){ q }() is: 123)
		Assert(function(_p = 0){ p }() is: 0)
		Assert(function(_p = 0){ p }() is: 0)
		Assert(function(_p = 0){ p }(123) is: 123)
		Assert(function(_p = 0){ p }(p: 123) is: 123)
		Assert(function(_q = 0){ q }() is: 123)
		}
	Test_dot_params()
		{
		c = class { New(.P) { } A() { .P } }
		i = c(123)
		Assert(i.A() is: 123)

		c = class { New(.p) { } A() { .p } }
		i = c(123)
		Assert(i.A() is: 123)
		}
	Test_combined()
		{
		_p = 123
		c = class { New(._P) { } A() { .P } }
		i = c()
		Assert(i.A() is: 123)
		i = c(456)
		Assert(i.A() is: 456)

		c = class { New(._p) { } A() { .p } }
		i = c()
		Assert(i.A() is: 123)
		i = c(456)
		Assert(i.A() is: 456)
		}
	Test_default_cant_be_identifier()
		{
		Assert({ "function (x = fred) { }".Compile() }
			throws: "parameter defaults must be constants")
		}
	}