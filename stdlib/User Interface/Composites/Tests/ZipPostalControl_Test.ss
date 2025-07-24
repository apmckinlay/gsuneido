// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(ZipPostalControl.ValidData?('1') is: false)
		Assert(ZipPostalControl.ValidData?('12345') is: true)
		Assert(ZipPostalControl.ValidData?('s1s 1a1') is: true)
		Assert(ZipPostalControl.ValidData?('') is: true)
		Assert(ZipPostalControl.ValidData?('', mandatory:) is: false)
		}
	}