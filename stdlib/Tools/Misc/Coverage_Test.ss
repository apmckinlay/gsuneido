// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		CoverageEnable(true)
		source = "function()
			{
			x = 0
			for (i = 0; i < 10; ++i)
				x += i
			if x is 0
				Print(x)
			return x
			}"
		f = source.Compile()
		f.StartCoverage(count:)
		f()
		cover = f.StopCoverage()
		Assert(cover is: #(21: 1, 31: 1, 61: 10, 72: 1, 87: 0, 100: 1))
		result = Coverage(source, cover)
		Assert(result like:
			"	function()
					{
			 1		x = 0
			 1		for (i = 0; i < 10; ++i)
			10			x += i
			 1		if x is 0
			 0			Print(x)
			 1		return x
					}")
		}
	Teardown()
		{
		CoverageEnable(false)
		super.Teardown()
		}
	}