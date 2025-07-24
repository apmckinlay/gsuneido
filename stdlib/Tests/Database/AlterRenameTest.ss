// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(one, two) key(one, two)')
		QueryOutput(table, #(one: 1, two: 2))
		Database("alter " $ table $ " rename two to too")
		QueryOutput(table, #(one: 11, too: 22))
		Assert(QueryAll(table),
			is: #((one: 1, too: 2), (one: 11, too: 22)))
		}
	}