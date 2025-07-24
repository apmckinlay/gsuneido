// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_scanNested()
		{
		scanNested = RazorCode.RazorCode_scanNested
		output = ""
		fn = {|s| output $= s }
		Assert(scanNested(0, fn, "(a + (b + 5)) more", '(')
			is: ' more')
		Assert(output is: 'a + (b + 5)')

		output = ""
		Assert(scanNested(0, fn, "{ { a } b } more", '{')
			is: ' more')
		Assert(output is: ' { a } b ')
		}
	}