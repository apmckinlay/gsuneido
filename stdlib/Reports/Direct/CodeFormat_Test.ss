// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_WrapDataLines()
		{
		wrap = CodeFormat.WrapDataLines
		_report = Mock()
		_report.When.GetCharWidth(1, '@mono', widthChar: '9').Return(10)
		Assert(wrap('', 1000, '@mono') is: '')
		Assert(wrap('a'.Repeat(10), 1000, '@mono') is: 'a'.Repeat(10))
		Assert(wrap('a'.Repeat(100), 1000, '@mono') is: '\r\n    ' $ 'a'.Repeat(100))
		Assert(wrap('abc '.Repeat(26), 1000, '@mono') is: 'abc '.Repeat(24).Trim() $
			'\r\n    abc abc ')
		}
	}