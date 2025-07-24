// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		a = new TestObserverString(quiet:)
		a.BeforeTest('test_name')
		a.BeforeMethod('test_method')
		a.Output('this is a test')
		Assert(a.Result is: 'this is a test\r\n')
		Assert(a.HasError?() is: false)
		}
	}