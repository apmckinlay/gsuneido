// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_emptykey()
		{
		table = .MakeTable('(a,b,c) key()')
		Assert(KeyQuery(table, #(a: 1, b: 2, c: 3)) is: table)
		}
	Test_singlekey()
		{
		table = .MakeTable('(a,b,c) key(b)')
		Assert(KeyQuery(table, #(a: 1, b: 2, c: 3)), is: table $ ' where b = 2')
		}
	Test_multikey()
		{
		table = .MakeTable('(a,b,c) key(b,a)')
		Assert(KeyQuery(table, #(a: 1, b: 2, c: 3)), is: table $ ' where b = 2 and a = 1')
		}
	}