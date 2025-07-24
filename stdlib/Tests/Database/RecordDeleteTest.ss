// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a, b, c) key(a)')
		for (i = 0; i < 5; ++i)
			QueryOutput(table, Object(a: i, b: i))
		QueryApply(table, update:)
			{ |x|
			x.Delete()
			}
		Assert(QueryEmpty?(table))
		}
	}