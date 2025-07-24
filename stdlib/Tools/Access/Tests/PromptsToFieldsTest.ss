// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		map = #(One: one, Two: two, Three: three)
		Assert(PromptsToFields(map, #()) is: #())
		Assert(PromptsToFields(map, #(Two)) is: #(two))
		Assert(PromptsToFields(map, #(Two, Three, One)) is: #(two, three, one))
		Assert(PromptsToFields(map, #(Two, 'Three ', One)) is: #(two, three, one))
		Assert(PromptsToFields(map, #(Two, 'jjj ', One)) is: #(two, one))
		}
	}