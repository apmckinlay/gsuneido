// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a, b, c, d, e, f, g) key(b)')
		ob = #(a: 1, b: 2, c: 3, d: 4, e: 5, f: 6, g: 7)
		QueryOutput(table, ob)
		Assert(Query1(table) is: ob)
		Database("alter " $ table $ " drop (a, c, e, g)")
		Database("alter " $ table $ " create (h, i, j, k)")
		Assert(Query1(table) is: #(b: 2, d: 4, f: 6))
		QueryDo('update ' $ table $ ' set b=1, d=2, f=3, h=4, i=5, j=6, k=7')
		Assert(Query1(table) is: #(b: 1, d: 2, f: 3, h: 4, i: 5, j: 6, k: 7))
		}
	Test_bug()
		{
		table = .MakeTable('(a,b,c) key(a)',
			r = [a: 1, b: 2, c: 3])
		Database("alter " $ table $ " create (d)")
		Database("alter " $ table $ " drop (d)")
		Database("alter " $ table $ " create (d)")
		Assert(Query1(table) is: r)

		QueryApply(table $ ' extend m=123', update:)
			{ |x|
			x.d = 4
			x.Update()
			}
		r.d = 4
		Assert(Query1(table) is: r)

		QueryDo('update ' $ table $ ' extend m=123 set d=5')
		r.d = 5
		Assert(Query1(table) is: r)

		QueryOutput(table, r = [a: 11, b: 22, c: 33, d: 44])
		Assert(QueryLast(table $ ' sort a') is: r)

		QueryDo('insert {a: 111, b: 222, c: 333, d: 444} into ' $ table)
		Assert(QueryLast(table $ ' sort a') is: [a: 111, b: 222, c: 333, d: 444])
		}
	}