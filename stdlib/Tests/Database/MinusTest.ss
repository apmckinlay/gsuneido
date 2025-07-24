// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table1 = .MakeTable('(a, b, c) key(a)')
		table2 = .MakeTable('(b, c, d) key(d)')
		Assert(QueryEmpty?(table1 $ ' minus ' $ table2), msg: 'one')
		QueryOutput(table1, #(b: 1, c: 2))
		QueryOutput(table2, #(b: 1, c: 2))
		Assert(QueryEmpty?(table1 $ ' minus ' $ table2), msg: 'two')
		QueryOutput(table1, #(a: 1, b: 2, c: 3))
		QueryOutput(table2, #(b: 2, c: 3, d: 4))
		Assert(Query1(table1 $ ' minus ' $ table2) is: #(a: 1, b: 2, c: 3))

		Assert(Query1('(' $ table1 $ ' minus ' $ table2 $ ') where a is 1')
			is: #(a: 1, b: 2, c: 3))
		Assert(QueryEmpty?('(' $ table1 $ ' minus ' $ table2 $ ') where a is ""'),
			msg: 'three')
		Assert(QueryEmpty?('(' $ table1 $ ' minus ' $ table2 $ ') where b is 9'),
			msg: 'four')
		Assert(QueryEmpty?('(' $ table1 $ ' minus ' $ table2 $ ') where a is 9'),
			msg: 'five')
		Assert({ Query1('(' $ table1 $ ' minus ' $ table2 $ ') where d is 1')}
			throws: 'nonexistent columns')
		}
	}
