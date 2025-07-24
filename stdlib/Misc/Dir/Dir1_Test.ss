// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Dir1("_non_existent_") is: false)
		tmp = .MakeFile("hello world")
		Assert(Dir1(tmp) is: tmp)
		Assert(Dir1(tmp, details:).size is: 11)
		.MakeFile("foo bar")
		Assert({ Dir1("tests_*") } throws: "Dir1 got more than one record")
		}
	}