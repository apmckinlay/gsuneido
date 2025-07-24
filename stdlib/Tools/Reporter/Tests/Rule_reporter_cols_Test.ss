// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rec = Record()
		rec.selectFields = class { PromptToField(prompt) { return prompt } }
		Assert(rec.reporter_cols is: #())

		rec.allcols = #(col1 col2 col3)
		Assert(rec.reporter_cols is: #(col1 col2 col3))

		rec.allcols = #(col1 col2 col3, col4_lower!)
		Assert(rec.reporter_cols is: #(col1 col2 col3))

		rec.nonsummarized_fields = #(col2)
		Assert(rec.reporter_cols is: #(col1 col3))
		}
	}