// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_rule()
		{
		Assert([].reporter_summarizeby_cols is: #())
		Assert([summarize_func_cols: ''].reporter_summarizeby_cols is: #())
		Assert([summarize_func_cols: #()].reporter_summarizeby_cols is: #())

		rpt = [reporter_cols: #(abc, efg), summarize_func_cols: #(abc)]
		Assert(rpt.reporter_summarizeby_cols is: #(efg))
		}
	}