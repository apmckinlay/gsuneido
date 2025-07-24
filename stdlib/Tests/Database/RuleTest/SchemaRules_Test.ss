// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		.MakeLibraryRecord(
			#(name: Rule_c, text: 'function () { .a + .b }'))
		table = .MakeTable('(a, b, C) key(a)',
			[a: 1, b: 10],
			[a: 2, b: 20],
			[a: 3, b: 30],
			[a: 4, b: 40])
		x = Query1(table, c: 22)
		Assert(x.c is: 22)
		Assert(x is: [a: 2, b: 20, c: 22])
		}
	Test_two()
		{
		.MakeLibraryRecord(
			#(name: Rule_corule, text: 'function () { true }'),
			#(name: Rule_cjrule, text: 'function () { true }'))
		co = .MakeTable('(kco, common, Corule) key(kco)',
			[kco: 1, common: 10])
		cj = .MakeTable('(kcj, common, Cjrule) key(kcj)',
			[kcj: 2, common: 20])

		query = co $
			' extend include = corule is true' $
			' where include is true'
		x = Query1(query)
		Assert(x is: [common: 10, kco: 1, include: true])

		query = '(' $
			' (' $ co $ ' extend r = corule)' $
			' union ' $
			' (' $ cj $ ' extend r = cjrule)' $
			' )' $
			' extend include = r is true' $
			' where include is true
			/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */'
		list = QueryAll(query).Sort!(By(#common))
		.eq(list, #(
			(common: 10, kco: 1, include: true, r: true, corule:),
			(common: 20, kcj: 2, include: true, r: true, cjrule:)))

		query = '(' $ co $
			' union ' $
			' (' $ cj $ ' extend corule = cjrule) )' $
			' extend include = corule is true' $
			' where include is true
			/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */'
		list = QueryAll(query).Sort!(By(#common))
		.eq(list, #(
			(common: 10, kco: 1, include: true),
			(common: 20, kcj: 2, include: true, corule: true)))
		}
	eq(rows, expected)
		{
		Assert(rows.Size() is: expected.Size())
		for i in ..rows.Size()
			{
			r = rows[i]
			e = expected[i]
			for m in e.Members()
				Assert(r[m] is: e[m], msg: "row " $ i)
			}
		}
	}
