// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(ConvertNumeric('') is: '')
		Assert(ConvertNumeric('3') is: 3)
		Assert(ConvertNumeric('0') is: 0)
		Assert(ConvertNumeric('03') is: '03')
		Assert(ConvertNumeric('Abba') is: 'Abba')
		Assert(ConvertNumeric('1000.01') is: 1000.01)
		Assert(ConvertNumeric('3e10') is: '3e10')
		Assert(ConvertNumeric(false) is: false)
		Assert(ConvertNumeric('99999999999999999999') is: '99999999999999999999')
		}
	}