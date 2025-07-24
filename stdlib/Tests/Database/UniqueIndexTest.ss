// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a,b) key(a) index unique(b)')
		QueryOutput(table, Record(a: 1))
		QueryOutput(table, Record(a: 2))

		QueryOutput(table, Record(a: 3, b: 1))
		QueryOutput(table, Record(a: 4, b: 2))
		Assert({QueryOutput(table, Record(a: 5, b: 2))} throws: "duplicate key")

		Assert({QueryDo('update ' $ table $ ' where b = 2 set b = 1')}
			throws: "duplicate key")
		}
	Test_two()
		{
		table = .MakeTable('(a,b,c) key(a) index unique(b, c)')
		QueryOutput(table, Record(a: 1))
		QueryOutput(table, Record(a: 2))

		QueryOutput(table, Record(a: 3, b: 1))
		QueryOutput(table, Record(a: 4, b: 1, c: 1))
		Assert({QueryOutput(table, Record(a: 5, b: 1))} throws: "duplicate key")
		Assert({QueryOutput(table, Record(a: 6, b: 1, c: 1))} throws: "duplicate key")
		Assert({QueryDo('update ' $ table $ ' where a = 3 set c = 1')}
			throws: "duplicate key")
		}
	Test_creation()
		{
		table = .MakeTable('(a,b) key(a) index unique(b)')
		Database("ensure " $ table $ " (c) index unique(c)")
		Database("alter " $ table $ " create (d) index unique(d)")
		QueryApply('indexes where columns isnt "a"', :table)
			{|x|
			Assert(x.key is: 'u')
			}
		}
	}