// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(StateProvNameToCode('') is: false)
		Assert(StateProvNameToCode('SK') is: false)
		Assert(StateProvNameToCode('Saskatchewan') is: 'SK')
		Assert(StateProvNameToCode('SASKATCHEWAN') is: 'SK')

		Assert(StateProvNameToCode('Florida') is: 'FL')
		Assert(StateProvNameToCode('FLORIDA') is: 'FL')
		}
	}