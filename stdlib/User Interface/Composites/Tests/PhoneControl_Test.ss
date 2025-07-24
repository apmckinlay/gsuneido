// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(PhoneControl.ValidData?('1') is: false)
		Assert(PhoneControl.ValidData?('555-5555') is: true)
		Assert(PhoneControl.ValidData?('555-555-5555') is: true)
		Assert(PhoneControl.ValidData?('555-555-5555x55') is: true)
		Assert(PhoneControl.ValidData?('1-555-555-5555') is: true)
		Assert(PhoneControl.ValidData?('15555555555') is: true)
		Assert(PhoneControl.ValidData?('2-555-555-5555') is: false)
		Assert(PhoneControl.ValidData?('1-555-555-5555x55') is: true)
		Assert(PhoneControl.ValidData?('2-555-555-5555x55') is: false)
		Assert(PhoneControl.ValidData?('52-555-5555-555x55') is: true)
		Assert(PhoneControl.ValidData?('1234567') is: true)
		Assert(PhoneControl.ValidData?('') is: true)
		Assert(PhoneControl.ValidData?('', mandatory:) is: false)
		}
	}