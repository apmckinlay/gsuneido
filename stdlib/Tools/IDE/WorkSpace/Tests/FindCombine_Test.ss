// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = FindCombine
			{
			FindCombine_stop?(@unused)
				{
				return false
				}
			FindCombine_evaluate(args)
				{
				return args
				}
			}
		Assert(cl() is: #())
		Assert(cl(#(2), #(3,4), #(1,2)) is: #(1, 2, 3, 4))
		}
	}