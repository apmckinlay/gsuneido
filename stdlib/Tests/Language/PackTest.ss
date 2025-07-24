// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// SuJsWebTest
	Test_one()
		{
		test = function (x)
			{
			x2 = Unpack(Pack(x))
			Assert(x2 is: x)
			Assert(Type(x2) is: Type(x))
			}
		test(true)
		test(false)
		test(123) // integer
		test(1e66)
		test("")
		test("hello\x00world")
		test(#20150429)
		test(Date())
		test(Date.Begin())
		test(Date.End())
		test(#())
		test(Object())
		test(Record())
		test(#(1, 2, a: 3, b: 4))
		test(#((1), a: (2)))

		x = Object(1, 2, a: 3, b: 4)
		x.Set_default(0)
		x.Set_readonly()
		test(x)
		Assert(x.y is: 0)
		Assert({ x.y = 123 } throws: "can't modify readonly objects")
		x2 = Unpack(Pack(x))
		x2.y = 123 // lost read-only
		Assert({ x2.z } throws: "member not found") // lost default

		Assert({ Pack(this) } throws: "can't pack")
		Assert({ Pack(Test) } throws: "can't pack")
		}
	}