// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.table = .MakeTable("(a,b,c) key(a,b)")
		for a in .. 4
			for b in .. 4
					QueryOutput(.table, [:a, :b, c: 123])
		.assertNoTempindex("sort a")
		.assertNoTempindex("sort a,b")
		.assertNoTempindex("where a=1 sort a,b")
		.assertNoTempindex("where b=1 sort a,b")
		.assertNoTempindex("where c=1 sort a,b")
		.assertNoTempindex("where a=1 and b=1 sort a,b")
		.assertTempindex("sort a,b,c")
		.assertTempindex("sort b")
		.assertNoTempindex("where a=1 sort b")
		.assertNoTempindex("where c=1 sort a,b,c")
		}
	assertTempindex(query)
		{
		strategy = QueryStrategy(.table $ ' ' $ query)
		Assert(strategy has: "tempindex", msg: strategy)
		}
	assertNoTempindex(query)
		{
		strategy = QueryStrategy(.table $ ' ' $ query)
		Assert(strategy hasnt: "tempindex", msg: strategy)
		}
	}
