// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table1 = .MakeTable('(a, b, c) key(a)')
		table2 = .MakeTable('(b, c, d) key(d)')
		Assert(QueryEmpty?(table1 $ ' intersect ' $ table2), msg: 'one')
		QueryOutput(table1, #(a: 1, b: 2, c: 3))
		QueryOutput(table1, #(a: 2, b: 3, c: 4))
		QueryOutput(table2, #(b: 2, c: 3, d: 4))
		QueryOutput(table2, #(b: 3, c: 4, d: 5))
		Assert(QueryEmpty?(table1 $ ' intersect ' $ table2), msg: 'two')
		QueryOutput(table1, #(b: 1, c: 2))
		QueryOutput(table2, #(b: 1, c: 2))
		Assert(Query1(table1 $ ' intersect ' $ table2) is: #(b: 1, c: 2))

		Assert(Query1('(' $ table1 $ ' intersect ' $ table2 $ ') where c is 2')
			is: #(b: 1, c: 2))
		Assert(QueryEmpty?('(' $ table1 $ ' intersect ' $ table2 $ ') where b is 9'),
			msg: 'three')
		Assert({ Query1('(' $ table1 $ ' intersect ' $ table2 $ ') where a is 1')}
			throws: 'nonexistent columns')
		Assert({ Query1('(' $ table1 $ ' intersect ' $ table2 $ ') where d is 1')}
			throws: 'nonexistent columns')
		}
	}
