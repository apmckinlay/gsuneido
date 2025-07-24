// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// SuJsWebTest
	Test_one()
		{
		x = 0
		Assert(Finally({ 123 }, { x = 2 }) is: 123)
		Assert(x is: 2)

		Assert({ Finally({ }, { throw "f" }) } throws: "f")

		Assert({ Finally({ throw "m" }, { }) } throws: "m")

		Assert({ Finally({ throw "e1" }, { throw "e2" }) } throws: "e1")

		x = 0
		Finally({  }, { ++x })
		Assert(x is: 1)
		try
			Finally({ throw "x" }, { ++x })
		catch (unused, "x")
			{}
		Assert(x is: 2)

		.n = 0
		.f()
		Assert(.n is: 1)

		Assert(.g() is: 123)
		Assert(.n is: 456)
		}
	f()
		{
		Finally({  }, { ++.n; return })
		}
	g()
		{
		Finally({ return 123 }, { .n = 456 })
		}

	Test_return_throw()
		{
		// SuJsWebTest Excluded - return throw is not implemented in SuJs
		unused = Finally(function() { return throw 'test1' }, { }) // doesn't throw

		Assert({ Finally(function() { return throw 'test1' }, { }) }
			throws: "test1")

		Assert({ Finally({ }, function() { return throw 'test2' }) }
			throws: "test2")

		Assert({ unused = Finally({ }, function() { return throw 'test3' }) }
			throws: "test3")

		Assert({ Finally({ }, function() { return throw 'test4' }) }
			throws: "test4")
		}
	}