// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(PostalControl.ValidData?('1') is: false)
		Assert(PostalControl.ValidData?('s1s 1a1') is: true)
		Assert(PostalControl.ValidData?('') is: true)
		Assert(PostalControl.ValidData?('', mandatory:) is: false)
		}
	}