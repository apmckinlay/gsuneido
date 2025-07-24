// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(HandleSetConcurrentError(#()) is: #())
		Assert(HandleSetConcurrentError(x = #(1, 2, a: 3, b: 4)) is: x)
		Assert(HandleSetConcurrentError([Md5()]) is: #("md5"))
		Assert(HandleSetConcurrentError([mem: Md5()]) is: #(mem: "md5"))
		Assert(HandleSetConcurrentError([a: 1, b: Object(mem: Md5()), c: [m: Md5()]])
			is: [a: 1, b: #(mem: "md5"), c: [m: "md5"]])
		}
	}
