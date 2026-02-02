// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(ZipPostalControl.ValidData?('1') is: false)
		Assert(ZipPostalControl.ValidData?('12345'))
		Assert(ZipPostalControl.ValidData?('s1s 1a1'))
		Assert(ZipPostalControl.ValidData?(''))
		Assert(ZipPostalControl.ValidData?('', mandatory:) is: false)
		}
	}