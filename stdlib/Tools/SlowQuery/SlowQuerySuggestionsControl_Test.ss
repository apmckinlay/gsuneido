// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_filterOperation()
		{
		f = SlowQuerySuggestionsControl.SlowQuerySuggestionsControl_filterOperation
		Assert(f('date') is 'greater than')
		Assert(f('num') is 'equals')
		Assert(f('number') is 'greater than')
		Assert(f('city') is 'equals')
		}
	}