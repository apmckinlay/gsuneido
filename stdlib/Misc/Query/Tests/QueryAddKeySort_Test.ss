// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable("(a, b, c) key()") // empty key
		Assert(QueryAddKeySort(table) is: table)
		Assert(QueryAddKeySort(table $ ' sort b') is: table $ ' sort b')

		table = .MakeTable("(a, b, c, d) key(a,b) key(c)")
		Assert(QueryAddKeySort(table)
			is: table $ ' sort c')
		Assert(QueryAddKeySort(table $ ' sort a,b')
			is: table $ ' sort a,b')
		Assert(QueryAddKeySort(table $ ' sort d,b,a')
			is: table $ ' sort d,b,a')
		Assert(QueryAddKeySort(table $ ' sort d') is: false)
		}
	}