// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Make()
		{
		c = TestData
			{
			N: 0
			TestData_output(@unused)
				{ ++.N }
			}
		td = new c
		td.TestData_output('test', 10, 10, 500)
		Assert(td.N is: 1)
		}
	Test_make_key()
		{
		mk = TestData.TestData_make_key
		Assert(mk(123, 5) isSize: 5)
		Assert(mk(123, 10) isSize: 10)
		Assert(mk(123, 20) isSize: 20)
		}
	}