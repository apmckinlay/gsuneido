// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		Assert(PostalControl.ValidData?('1') is: false)
		Assert(PostalControl.ValidData?('s1s 1a1'))
		Assert(PostalControl.ValidData?(''))
		Assert(PostalControl.ValidData?('', mandatory:) is: false)
		}
	}