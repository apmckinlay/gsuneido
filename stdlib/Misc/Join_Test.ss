// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Join('', 'one', 'two', 'three') is: 'onetwothree')
		Assert(Join(' ', 'one', 'two', 'three') is: 'one two three')
		Assert(Join('==', 'one', 'two', 'three') is: 'one==two==three')
		}
	}