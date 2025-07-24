// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(ZipControl.ValidData?('1') is: false)
		Assert(ZipControl.ValidData?('12345') is: true)
		Assert(ZipControl.ValidData?('') is: true)
		Assert(ZipControl.ValidData?('', mandatory:) is: false)
		}
	}