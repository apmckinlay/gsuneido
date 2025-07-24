// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_empty_table_bug()
		{
		table1 = .MakeTable("(a) key(a)")
		table2 = .MakeTable("(a) key(a)",
			[a: 1], [a: 2], [a: 3])
		WithQuery(table1 $ ' leftjoin by(a) ' $ table2 $ ' summarize count')
			{ |q /*unused*/| }
		}

	Test_noResults()
		{
		hdr = .MakeTable('(num) key(num)', [num: 0])
		lin = .MakeTable('(num, linnum, amount, date, string) key(linnum)',
			lin1 = [num: 0, linnum: 1, amount: 2, date: #19000101, string: 'string'],
			[num: 0, linnum: 2, amount: ""])

		Assert(Query1(lin $ ' where amount isnt ""')
			is: lin1 msg: 1)
		Assert(Query1(hdr $ ' join by(num) ' $ lin $ ' where amount isnt ""')
			is: lin1, msg: 2)
		Assert(Query1(hdr $ ' leftjoin by(num) ' $ lin $ ' where date isnt ""')
			is: lin1, msg: 3)
		Assert(Query1(hdr $ ' leftjoin by(num) ' $ lin $ ' where string isnt ""')
			is: lin1, msg: 4)
		Assert(Query1(hdr $ ' leftjoin by(num) ' $ lin $ ' where amount isnt ""')
			is: lin1, msg: 5)
		}
	}
