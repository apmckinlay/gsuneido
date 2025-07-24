// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// NOTE: this test uses a lot of database space
Test
	{
	Test_split_over_summarize()
		{
		table = .MakeTable('(a,b,c) key(a)', [a: 1, b: 2, c: 3])
		query = table $ '
			summarize a, total c
			where a > 1 and total_c > 100'
		Assert(QueryStrategy(query)
			is: table $ "^(a) where a > 1 summarize-seq a, total c where total_c > 100")
		}
	Test_split_over_extend()
		{
		table = .MakeTable('(a,b,c) key(a)')
		query = table $ '
			extend rule, d = b + c
			where b > 4 and d > 5 and rule < 3'
		Assert(QueryStrategy(query)
			is: table $ "^(a) where b > 4 and b + c > 5 " $
				"extend rule, d = b + c where rule < 3")
		}

	Test_iselects()
		{
		table = .MakeTable('(a,b,c) key(a,b,c)')
		for a in ..3
			for b in ..3
				for c in ..3
					QueryOutput(table, [:a, :b, :c])
		query = table $ ' where a = 1 and c = 1'
		QueryApply(query)
			{ |x|
			Assert(x.a is: 1)
			Assert(x.c is: 1)
			}
		}
	Test_fixed()
		{
		table = .MakeTable('(a,b,c) key(a)')
		query = table $ ' where b = 1 sort a,b'
		Assert(QueryStrategy(query) hasnt: 'tempindex')
		query = table $ ' where b = 1 sort b,a'
		Assert(QueryStrategy(query) hasnt: 'tempindex')
		}

//	Test_select_conflict_with_fixed() // what conflict ???
//		{
//		bizpartners = .MakeTable('(biznum) key(biznum)',
//			[biznum: 1], [biznum: 2])
//		orders = .MakeTable('(ordernum, biznum, invoice) key(ordernum) index(invoice)',
//			[ordernum: 1, biznum: 1, invoice: 1],
//			[ordernum: 2, biznum: 2, invoice: 2],
//			[ordernum: 3, biznum: 1, invoice: 1],
//			[ordernum: 4, biznum: 2, invoice: 2])
//		query = bizpartners $ ' join by(biznum) ' $
//			'(' $ orders $ ' where biznum = 1 and invoice = 1)'
//		Assert(QueryStrategy(query).Tr('()')
//			is: orders $ "^invoice WHERE biznum is 1 and invoice is 1 " $
//				"JOIN n:1 on biznum " $ bizpartners $ "^biznum")
//		WithQuery(query)
//			{ |q|
//			Assert(q.Next() is: [invoice: 1, ordernum: 1, biznum: 1])
//			Assert(q.Next() is: [invoice: 1, ordernum: 3, biznum: 1])
//			Assert(q.Next() is: false)
//			}
//		}
	}
