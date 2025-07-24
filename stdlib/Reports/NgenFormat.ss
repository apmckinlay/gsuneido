// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// generates data for testing
Generator
	{
	New(.n = 100)
		{
		.i = 0
		}
	Next()
		{
		if (++.i <= .n)
			return _report.Construct(Object('Text', .i))
		else
			return false
		}
	}