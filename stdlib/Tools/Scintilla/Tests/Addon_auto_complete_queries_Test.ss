// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		mock = Mock(Addon_auto_complete_queries)
		mock.When.IdleAfterChange().CallThrough()
		mock.When.AutoComplete([anyArgs:]).CallThrough()
		mock.When.AutoShow([anyArgs:]).Do({ })

		table1 = .MakeTable('(field_a1, field_b1, special) key(field_a1)')
		table2 = .MakeTable('(field_a2, field_b2, unique) key(field_a2)')
		table3 = .MakeTable('(field_a3, field_b3, other) key(field_a3)')
		mock.When.collectTables().Return(tables = Object(table1, table2, table3))

		Assert(mock.Addon_auto_complete_queries_tables is: #())
		Assert(mock.Addon_auto_complete_queries_fields is: #())

		mock.When.Get().Return(table1)
		mock.IdleAfterChange()
		Assert(mock.Addon_auto_complete_queries_tables is: tables)
		Assert(mock.Addon_auto_complete_queries_fields
			is: #('field_a1', 'field_b1', 'special'))

		mock.AutoComplete('f')
		mock.Verify.Never().matchingWords([anyArgs:])
		mock.Verify.Never().AutoShow([anyArgs:])

		mock.AutoComplete('fi')
		mock.Verify.matchingWords('fi')
		mock.Verify.AutoShow('fi', #('field_a1', 'field_b1'))

		mock.AutoComplete('sp')
		mock.Verify.matchingWords('sp')
		mock.Verify.AutoShow('sp', #('special'))

		mock.When.Get().Return(table2)
		mock.IdleAfterChange()
		Assert(mock.Addon_auto_complete_queries_tables is: tables)
		Assert(mock.Addon_auto_complete_queries_fields
			is: #('field_a2', 'field_b2', 'unique'))

		mock.AutoComplete('sp')
		mock.Verify.Times(2).matchingWords('sp')
		mock.Verify.AutoShow('sp', #())

		mock.AutoComplete('field_')
		mock.Verify.matchingWords('field_')
		mock.Verify.AutoShow('field_', #('field_a2', 'field_b2'))

		mock.AutoComplete('field_a')
		mock.Verify.matchingWords('field_a')
		mock.Verify.AutoShow('field_a', #('field_a2'))
		}
	}