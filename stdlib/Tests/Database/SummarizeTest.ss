// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(rep,date,amount) key(rep,date)')
		QueryOutput(table, #(rep: 1, date: 1, amount: 11))
		QueryOutput(table, #(rep: 1, date: 4, amount: 12))
		QueryOutput(table, #(rep: 2, date: 2, amount: 21))
		QueryOutput(table, #(rep: 2, date: 3, amount: 23))
		x = Query1(table $ " summarize count")
		Assert(x is: #(count: 4))

		x = Query1(table $ " summarize total amount")
		Assert(x is: #(total_amount: 67))

		x = Query1(table $ " summarize average amount")
		Assert(x is: #(average_amount: 16.75))

		x = Query1(table $ " summarize min date")
		Assert(x.min_date is: 1)

		x = Query1(table $ " summarize max date")
		Assert(x.max_date is: 4)

		x = Query1(table $ " summarize
			count, total amount, average amount, min date, max date")
		Assert(x
			is: #(count: 4, total_amount: 67, average_amount: 16.75,
				min_date: 1, max_date: 4))

		byrep = #(
			1: (rep: 1, count: 2, min_date: 1, max_date: 4, total_amount: 23,
				average_amount: 11.5)
			2: (rep: 2, count: 2, min_date: 2, max_date: 3, total_amount: 44,
				average_amount: 22)
			)
		QueryApply(table $ " summarize rep, count, min date, max date,
			total amount, average amount")
			{|x|
			Assert(x is: byrep[x.rep])
			}

		strategy = QueryStrategy(table $ " where rep = 1
			summarize date, total amount")
		Assert(strategy has: "^(rep,date)")
		}
	Test_fixed()
		{
		table = .MakeTable("(a,b,c,d) key(a,b) index(a) index(b)")
		for i in .. 10
			QueryOutput(table, Record(a: i.Even?(), b: i))
		query = table $ ' where a = true summarize b, min d sort b'
		Assert(QueryStrategy(query) has: table $ '^(a,b)')
		}
	Test_map()
		{
		cases = #(
			'columns summarize column, count',
			'indexes summarize table, key, count',
			)
		for c in cases
			{
			Assert(QueryStrategy(c) has: 'summarize-map')
			Assert(QueryStrategy(c) hasnt: 'tempindex')
			}
		}
	Test_idx()
		{
		q = 'stdlib summarize min num' // key
		Assert(QueryStrategy(q) has: 'stdlib^(num) summarize-idx')
		Assert(Query1(q).name isnt: '') // whole record

		q = 'stdlib summarize min num sort lib_modified' // ignore sort (only 1 record)
		strategy = QueryStrategy(q)
		Assert(strategy has: 'stdlib^(num) summarize-idx')
		Assert(strategy hasnt: 'tempindex')
		Assert(Query1(q).name isnt: '') // whole record

		q = 'stdlib summarize max parent' // index
		Assert(QueryStrategy(q) has: 'summarize-idx')
		Assert(Query1(q).Size() is: 1) // no record

		q = 'stdlib where num < 10 summarize max name'
		Assert(QueryStrategy(q) has: '-seq') // better to use num index
		q = 'stdlib where num < 1000 summarize max name'
		Assert(QueryStrategy(q) has: '-idx')
		}
	}
