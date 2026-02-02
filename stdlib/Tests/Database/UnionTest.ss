// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.table1 = .MakeTable('(num, text) key(num)')
		.table2 = .MakeTable('(num, text) key(num)')
		.table3 = .MakeTable('(num, text) key(num)')
		}

	Test_keys()
		{
		Assert(QueryKeys(.table1) is: #("num"))
		Assert(QueryKeys(.union(.table1, .table2)) is: #("num,text"))
		Assert(QueryKeys(.union(.union(.table1, .table2), .table3))
			is: #("num,text"))
		Assert(QueryKeys(.union(
			.table1 $ ' extend table = "one"',
			.table2 $ ' extend table = "two"')) is: #("num,table"))
		}

	Test_rewind()
		{
		table1 = .MakeTable('(one, two) key(one)')
		table2 = .MakeTable('(one, two) key(two)')
		for i in ..10
			{
			QueryOutput(table1, Object(one: i, two: i))
			QueryOutput(table2, Object(one: -i, two: -i))
			}
		query = .union(table1, table2)
		Transaction(read:)
			{ |t|
			q = t.Query(query)
			x = q.Prev()

			q = t.Query(query)
			q.Next()
			q.Rewind()
			y = q.Prev()
			}
		Assert(y is: x)
		}

	Test_union_where_bug()
		{
		QueryApply(.union('tables', 'columns') $ ' where table < "m"
			where table.Has?("x")')
			{ |x|
			Assert(x.table has: 'x')
			}
		}

	Test_different_columns()
		{
		table1 = .MakeTable('(a, b, c) key(a)')
		table2 = .MakeTable('(b, c, d) key(d)')
		Assert(QueryEmpty?(.union(table1, table2)), msg: 'table1 and table2')

		QueryOutput(table1, #(b: 1, c: 2))
		QueryOutput(table2, #(b: 1, c: 2))
		Assert(Query1(.union(table1, table2)) is: #(b: 1, c: 2))
		QueryOutput(table1, #(a: 1, b: 2, c: 3))
		QueryOutput(table2, #(b: 2, c: 3, d: 4))
		Assert(QueryAll(.union(table1, table2))
			is: #((a: 1, b: 2, c: 3), (b: 1, c: 2), (b: 2, c: 3, d: 4)))

		Assert(QueryAll(.union(table1, table2) $ ' where a is 1')
			is: #((a: 1, b: 2, c: 3)))
		Assert(QueryAll(.union(table1, table2) $ ' where a is ""')
			equalsSet: #((b: 1, c: 2), (b: 2, c: 3, d: 4)))
		Assert(QueryEmpty?(.union(table1, table2), b: 9), msg: 'b is 9')
		Assert(QueryEmpty?(.union(table1, table2), a: 9), msg: 'a is 9')
		Assert(QueryEmpty?(.union(table1, table2), d: 9), msg: 'd is 9')
		}

	Test_union_disjoint_bug()
		{
		tbl = .MakeTable("(k,x) key(k)")
		for k in .. 20
			QueryOutput(tbl, [:k, x: -k])
		query = .union(
			tbl $ ' where k = 3 remove k',
			tbl $ ' where k = 3 remove k')
		Assert(QueryStrategy(query) hasnt: 'disjoint')
		Assert(QueryCount(query) is: 1)

		query = .union(
			.union(tbl $ ' where k = 1', tbl $ ' where k = 2'),
			.union(tbl $ ' where k = 3', tbl $ ' where k = 4'))
		Assert(QueryStrategy(query) has: 'disjoint')
		Assert(QueryCount(query) is: 4)

		query = .union(
			.union(
				tbl $ ' where k = 3',
				tbl $ ' where k = 4'),
			tbl $ ' where k = 4')
		Assert(QueryStrategy(query) has: 'disjoint')
		Assert(QueryCount(query) is: 2)

		query = .union(
			.union(
				tbl $ ' where k = 4',
				tbl $ ' where k = 3'),
			.union(
				tbl $ ' where k = 4',
				tbl $ ' where k = 3'))
		Assert(QueryStrategy(query) =~ 'disjoint.*union-merge.*?disjoint')
		Assert(QueryCount(query) is: 2)

		query = .union(
			.union(tbl $ ' where k = 1', tbl $ ' where k = 2') $ 'remove k',
			.union(tbl $ ' where k = 3', tbl $ ' where k = 4') $ 'remove k') $
			'/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */'
		Assert(QueryStrategy(query) =~
			'disjoint.*union-(lookup|merge).*?disjoint')
		Assert(QueryCount(query) is: 4)

		query = .union(
			.union(
				tbl $ ' where k = 4 remove k',
				tbl $ ' where k = 3 remove k'),
			.union(
				tbl $ ' where k = 4 remove k',
				tbl $ ' where k = 3 remove k'))
		Assert(QueryCount(query) is: 2)
		}
	union(query1, query2)
		{
		return '((' $ query1 $ ') union (' $ query2 $ '))
			/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */'
		}

	Test_disjoint_keys()
		{
		table1 = .MakeTable('(a, b, c, x) key(a) key(x)')
		table2 = .MakeTable('(b, c, d, y) key(d) key(c,d)')
		Assert(QueryKeys(.union(table1, table2)) is: #("a,b,c,x,d,y"))
		disjoint = .union(table1 $ " where b = 1", table2 $ " where b = 2")
		Assert(QueryKeys(disjoint) is: #("a,d,b", "x,d,b"))
		}

	Test_keys_bug()
		{
		table = .MakeTable('(a,b,c) key(a,b) key(b,a)')
		Assert(QueryKeys(table) is: #("a,b"))

		q = "(" $ table $ " extend d=1) union (" $ table $ " extend d=2)"
		Assert(QueryKeys(q) is: #("a,b,d"))

		q = "(" $ table $ " extend d=1)
			union (" $ table $ " extend d=2)
			union (" $ table $ " extend d=3)"
		Assert(QueryKeys(q) is: #("a,b,d"))
		}
	Test_rule_on_one_source()
		{
		tbl1 = .MakeTable("(k, test_simple_default) key(k)",
			[k: 1, test_simple_default: 'aaa'],
			[k: 2, test_simple_default: 'zzz'])
		tbl2 = .MakeTable("(k, Test_simple_default) key(k)",
			[k: 3])
		list = Object()
		// intentionally not QueryList or QueryAll
		QueryApply(.union(tbl1, tbl2) $ " sort test_simple_default")
			{|x|
			list.Add(x.test_simple_default)
			}
		Assert(list is: list.Copy().Sort!())
		}
	}
