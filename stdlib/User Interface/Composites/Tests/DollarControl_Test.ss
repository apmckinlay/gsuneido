// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_convert()
		{
		convert = DollarControl.DollarControl_convert
		Assert(convert('$123,456') is: '123456')
		Assert(convert('(123,456)') is: '-123456')
		Assert(convert('123456.79') is: '123456.79')
		Assert(convert('$1,234,567.89') is: '1234567.89')
		Assert(convert('abc') is: 'abc')
		}
	}