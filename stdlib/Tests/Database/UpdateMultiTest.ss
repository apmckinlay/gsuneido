// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_update_update()
		{
		table = .MakeTable('(a, b, c) key(a)')
		Transaction(update:)
			{ |t|
			q = t.Query(table)
			q.Output(#(a: 1, b: 2, c: 3))
			x = q.Next()
			x.b = 22
			x.Update()
			x.a = 11
			x.Update()
			x.c = 33
			x.Update()
			}
		Assert(Query1(table $ ' where a = 11') is: #(a: 11, b: 22, c: 33))
		}

	Test_update_delete()
		{
		tbl = .MakeTable('(a, b) key (a)', [a: 1, b: 2])
		QueryApply1(tbl)
			{|x|
			x.b = 22
			x.Update()
			x.Delete()
			}
		Assert(QueryEmpty?(tbl), msg: 'not empty')
		}
	}
