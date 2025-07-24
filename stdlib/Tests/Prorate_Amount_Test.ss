// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Prorate_Amount(#(), 100) is: #())
		Assert(Prorate_Amount(#(10 10 10), 100) is: #(33.33 33.34 33.33))
		Assert(Prorate_Amount(#(-10 -10 -10), -100) is: #(-33.33 -33.34 -33.33))
		Assert(Prorate_Amount(#(12 13 14), 500) is: #(153.85 166.66 179.49))
		Assert(Prorate_Amount(#(244.56 123.24), 1850) is: #(1230.11 619.89))
		Assert(Prorate_Amount(#(amount1: 12, amount2: 13, amount3: 14), 500)
			is: #(amount3: 179.49, amount1: 153.85, amount2: 166.66))
		Assert(Prorate_Amount(#(test_one:244.56, test_two: 123.24), 1850)
			is: #(test_two: 619.89, test_one: 1230.11))
		Assert(Prorate_Amount(#(0 0 0), 100) is: #(33.33 33.34 33.33))
		Assert(Prorate_Amount(#(-10 10), 100) is: #(50, 50))
		Assert(Prorate_Amount(#(-10 10 10), 100) is: #(33.33 33.34 33.33))
		Assert(Prorate_Amount(#(10 10 10), 0) is: #(0, 0, 0))
		Assert(Prorate_Amount(#(10 0 10), 100) is: #(50 0 50))
		Assert(Prorate_Amount(#(10 0 10), -100) is: #(-50 0 -50))
		Assert(Prorate_Amount(#(-10 10 10), -100) is: #(-33.33 -33.34 -33.33))
		Assert(Prorate_Amount(#(a: 0, b: 0), 100) is: #(a: 50, b: 50))
		}
	}