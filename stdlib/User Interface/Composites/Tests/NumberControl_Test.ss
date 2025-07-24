// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData()
		{
		Assert(NumberControl.ValidData?("") is: true)
		Assert(NumberControl.ValidData?("", mandatory:) is: false)
		Assert(NumberControl.ValidData?('5') is: true)
		Assert(NumberControl.ValidData?('fred') is: false)
		Assert(NumberControl.ValidData?('1e99') is: true)
		Assert(NumberControl.ValidData?('1e999') is: false)
		Assert(NumberControl.ValidData?(0xB5.Chr()) is: false)
		// for the following 3 .Number?() returns true, but they cannot be converted
		// using Number()
		Assert(NumberControl.ValidData?(0xB2.Chr()) is: false)
		Assert(NumberControl.ValidData?(0xB3.Chr()) is: false)
		Assert(NumberControl.ValidData?(0xB9.Chr()) is: false)
		}

	Test_isSimpleMathExpression?()
		{
		m = NumberControl.NumberControl_isSimpleMathExpression?
		Assert(m('') is: false)
		Assert(m('test') is: false)
		Assert(m('1 + 2') is: true)
		Assert(m('1 - 2') is: true)
		Assert(m('1 * 2') is: true)
		Assert(m('1 / 2') is: true)
		Assert(m('1e2 / 2') is: true)
		Assert(m('6785 / 11.223') is: true)
		Assert(m('(3+6785) / 11.223') is: true)
		Assert(m('1/2') is: true)
		Assert(m('1^2') is: false)
		Assert(m('1+1;Print("test")') is: false)
		Assert(m('Print("test");1+1') is: false)
		Assert(m('Print("1+1")') is: false)
		Assert(m('123') is: false)
		Assert(m('0123') is: false)
		Assert(m('123.123') is: false)
		Assert(m('0123.123') is: false)
		Assert(m('0x132d') is: false)
		}
	}
