// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	NoRet()
		{ }

	Test_passthrough()
		{
		f = function ()
			{ ReturnTest.NoRet() }
		f()

		f = function ()
			{ (ReturnTest.NoRet()) }
		f()

		f = function ()
			{ return ReturnTest.NoRet() }
		f()

		f = function ()
			{ return (ReturnTest.NoRet()) }
		f()

		f = function ()
			{ unused = ReturnTest.NoRet() }
		Assert({ f() } throws: "no return value")

		f = function ()
			{ (unused = ReturnTest.NoRet()) }
		Assert({ f() } throws: "no return value")

		f = function (cond)
			{ cond ? ReturnTest.NoRet() : ReturnTest.NoRet() }
		f(true)
		f(false)

		f = function (cond)
			{ (cond ? ReturnTest.NoRet() : ReturnTest.NoRet()) }
		f(true)
		f(false)

		f = function (cond)
			{ return cond ? ReturnTest.NoRet() : ReturnTest.NoRet() }
		f(true)
		f(false)

		f = function (cond)
			{ return (cond ? ReturnTest.NoRet() : ReturnTest.NoRet()) }
		f(true)
		f(false)

		f = function (cond)
			{ unused = cond ? ReturnTest.NoRet() : ReturnTest.NoRet() }
		Assert({ f(true) } throws: "no return value")
		Assert({ f(false) } throws: "no return value")

		f = function (cond)
			{ unused = (cond ? ReturnTest.NoRet() : ReturnTest.NoRet()) }
		Assert({ f(true) } throws: "no return value")
		Assert({ f(false) } throws: "no return value")

		f = function (cond)
			{ cond ? ReturnTest.NoRet() : ReturnTest.NoRet(); 123 }
		f(true)
		f(false)

		f = function (cond)
			{ (cond ? ReturnTest.NoRet() : ReturnTest.NoRet()); 123 }
		f(true)
		f(false)

		f = function () { }; b = { return f() }; b() // block return
		}

	Test_return_throw()
		{
		// return throw requires a result
		Assert({ "function() { return throw }".Compile() } throws: "syntax error")

		.rt("") // should not throw
		.rt(true) // should not throw
		unused = .rt("error") // should not throw
		Assert({ .rt("error") } throws: "error")

		unused = .passthru1() // should not throw
		unused = .passthru2() // should not throw
		Assert({ .passthru1() } throws: "error")
		Assert({ .passthru2() } throws: "error")
		Assert({ .discard() } throws: "error")
		}

	rt(x)
		{
		return throw x
		}

	passthru1()
		{
		.rt("error")
		}

	passthru2()
		{
		return .rt("error")
		}

	discard()
		{
		.rt("error")
		return
		}
	}