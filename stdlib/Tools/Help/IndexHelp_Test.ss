// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_skipClasses()
		{
		skipClasses = IndexHelp.IndexHelp_skipClasses
		Assert(skipClasses("") is: #())
		Assert(skipClasses(".foo { display: none; }\n.bar { display: none; }")
			is: #(foo:, bar:))
		}
	Test_Foreach_record()
		{
		c = IndexHelp
			{
			IndexHelp_queryApply(helpbook/*unused*/, extraWhere/*unused*/, block)
				{
				block(.X)
				}
			IndexHelp_htmlWrapPrefix(helpbook/*unused*/)
				{
				return ""
				}
			IndexHelp_skipClasses(styles/*unused*/)
				{
				return #(foo:, bar:)
				}
			}
		test =
			{|input, expected|
			ih = new c
			ih.X = [name: "Title", text: input]
			result = ""
			ih.Foreach_record("book")
				{|x/*unused*/, text|
				result $= text
				}
			Assert(result like: expected)
			}
		test('<p>Now is the time for all</p>'
			'Now is the time for all')
		test('<p class="other">Now is the time for all</p>'
			'Now is the time for all')
		test('<p class="foo">Now is the time for all</p>'
			'')
		}
	}