// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Make_where()
		{
		table = .MakeTable('(a,b,c) key(a,c)')
		exp = ExplorerListModelLinked(table, #(a, c))
		exp.SetBaseQuery(table)
		Assert(exp.Make_where([]) is: table $ ' where a is "" and c is ""')
		Assert(exp.Make_where([a: 1 b: 2, c: 'abc'])
			is: table $ ' where a is 1 and c is "abc"')
		}
	}