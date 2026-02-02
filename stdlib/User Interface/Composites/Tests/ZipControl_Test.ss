// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(ZipControl.ValidData?('1') is: false)
		Assert(ZipControl.ValidData?('12345'))
		Assert(ZipControl.ValidData?(''))
		Assert(ZipControl.ValidData?('', mandatory:) is: false)
		}
	}