// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// for _lower! columns
Test
	{
	Test_plain_record()
		{
		Assert([].a_lower! is: '') // missing
		Assert([a: #()].a_lower! is: #()) // non-string
		Assert([a: 123].a_lower! is: 123) // non-string
		Assert([a: "AbcD"].a_lower! is: 'abcd')
		}
	Test_main()
		{
		Assert({ .MakeTable('(a, b_lower!) key(a)') }
			throws: 'nonexistent column: b')

		// unindexed
		table = .MakeTable('(a, b, b_lower!) key(a)')
		Assert(Schema(table).Tr('\n') like: table $ ' (a, b, b_lower!) key (a)')
		.test(table)

		// indexed
		table = .MakeTable('(a, b, b_lower!) key(b_lower!)')
		idx = Query1('indexes', :table)
		Assert(idx.columns is: "b_lower!")
		Assert(Schema(table).Tr('\n') like: table $ ' (a, b, b_lower!) key (b_lower!)')
		.test(table, ' indexed')
		}

	test(tbl, msg = "")
		{
		Assert(QueryColumns(tbl) equalsSet: #(a, b, b_lower!),
			msg: 'Columns' $ msg)
		WithQuery(tbl) {|q|
			Assert(q.RuleColumns() is: #(b_lower!), msg: 'RuleColumns' $ msg) }

		QueryOutput(tbl, [b: 'FooBar'])
		Assert({ QueryOutput(tbl, [b: 'FOOBAR']) } throws: "duplicate")

		Assert(Query1(tbl, b: 'foobar') is: false)
		Assert(Query1(tbl, b_lower!: 'foobar') is: [b: 'FooBar'])
		Assert(Query1(tbl $ ' where b_lower!.Has?("oba")') is: [b: 'FooBar'])

		QueryOutput(tbl, [a: 2, b: 123]) // non-string
		}
	}
