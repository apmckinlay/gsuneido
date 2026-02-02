// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Regex?(true) is: false)
		Assert(Regex?(''))
		Assert(Regex?('abc'))
		Assert(Regex?('(abc)'))

		Assert(Regex?('(abc') is: false)
		Assert(Regex?('[abc') is: false)
		Assert(Regex?('abc)') is: false)
		}
	}