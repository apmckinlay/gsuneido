// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		.check("Test_test1", true)
		.check("Test_test1!", true)
		.check("Test_test1?", true)
		.check("Test_test1!_test", true)
		.check("Test_test1?_test", true)
		.check("Test.css", true)
		.check("test.js", true)

		.check("test_test1", false)
		.check("Test_test1?_test?", false)
		.check("test?.js", false)
		}
	check(testStr, expected)
		{
		Assert(testStr =~ '^' $ GlobalRegExForGoTo $ '$' is: expected)
		}
	}