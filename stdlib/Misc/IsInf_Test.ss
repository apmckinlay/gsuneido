// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(IsInf?(-1) is: false)
		Assert(IsInf?(0) is: false)
		Assert(IsInf?(1) is: false)
		Assert(IsInf?("inf") is: false)

		Assert(IsInf?(1/0))
		Assert(IsInf?(-1/0))
		}
	}