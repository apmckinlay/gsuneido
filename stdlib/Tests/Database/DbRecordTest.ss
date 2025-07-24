// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Update_with_separate_record()
		{
		tbl = .MakeTable("(a, b) key(a)", [a: 1, b: 2])
		Transaction(update:)
			{|t|
			x = t.QueryFirst(tbl $ ' sort a')
			x.Update([a: 3, b: 4])
			}
		Assert(QueryFirst(tbl $ ' sort a') is: [a: 3, b: 4])
		}
	}