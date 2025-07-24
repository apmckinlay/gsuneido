// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(FindByTokenScan("") is: #())
		Assert(FindByTokenScan(" /* */ \r\n") is: #())
		Assert(FindByTokenScan(". Sort (") is: #('.', 'Sort', '('))
		Assert(FindByTokenScan("'foo'") is: #('"foo"'))
//		Assert(FindByTokenScan("#foo") is: #('"foo"'))
		}
	}