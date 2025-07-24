// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_noexistent_table()
		{
		nonExistingTable = "Schema_Test_" $ Display(Timestamp())
		Assert({ Schema(nonExistingTable) }
			throws: "Schema: non-existent table: " $ nonExistingTable)
		}

	Test_simple_table()
		{
		table = .MakeTable('(a,b,c) index (b) key(a)')
		schema = Schema(table)
		Assert(schema has: '(a, b, c)')
		Assert(schema has: 'index (b)')
		Assert(schema has: 'key (a)')
		}

	Test_foreign_key_cascade()
		{
		table = .MakeTable('(a,b,c) index (b) key(a)')
		table2 = .MakeTable('(d,e,f,a)
			index (a) in ' $ table $ ' cascade
			key(d)')

		schema = Schema(table2)
		Assert(schema has: '(a, d, e, f)')
		Assert(schema has: 'index (a) in ' $ table $ ' cascade')
		Assert(schema has: 'key (d)')
		}

	Test_rule_field()
		{
		table = .MakeTable('(a,b,c,Schema_test_field) index (b) key(a)')
		schema = Schema(table)
		Assert(schema has: '(a, b, c, Schema_test_field)')
		}
	}