// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_args()
		{
		f = function (unused) { }

		Assert({ f() } throws: "missing argument")
		f(1)
		Assert({ f(1, 2) } throws: "too many arguments")

		Assert({ f(@#()) } throws: "missing argument")
		f(@#(1))
		Assert({ f(@#(1, 2)) } throws: "too many arguments")

		Assert({ f(@+1 #(1)) } throws: "missing argument")
		f(@+1 #(1, 2))
		Assert({ f(@+1 #(1, 2, 3)) } throws: "too many arguments")
		}
	}