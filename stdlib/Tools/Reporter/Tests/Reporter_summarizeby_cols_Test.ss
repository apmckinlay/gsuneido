// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rec = Record(reporter_cols: #(one two three), summarize_func_cols: #(three))
		Assert(rec.reporter_summarizeby_cols is: #(one two))
		}
	}