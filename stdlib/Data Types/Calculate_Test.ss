// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_PercentOf()
		{
		Assert(Calculate.PercentOf(101, '33', 1) is: 33.3)
		Assert(Calculate.PercentOf(101, '33', 2) is: 33.33)
		Assert(Calculate.PercentOf(101, '33') is: 33.33)
		Assert(Calculate.PercentOf(100, 50) is: 50)
		Assert(Calculate.PercentOf(100, 50) is: 50)
		Assert(Calculate.PercentOf(5, '25') is: 1.25)
		}

	Test_ConvertToDollars()
		{
		Assert(Calculate.ConvertToDollars(5) is: .05)
		Assert(Calculate.ConvertToDollars(50) is: .50)
		Assert(Calculate.ConvertToDollars('50') is: .50)
		Assert(Calculate.ConvertToDollars('5') is: .05)
		Assert(Calculate.ConvertToDollars('150') is: 1.50)
		Assert(Calculate.ConvertToDollars('32456') is: 324.56)
		}
	}