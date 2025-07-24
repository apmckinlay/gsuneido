// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	snapshotTables: SnapshotTables
		{
		SnapshotTables_logTable: 	'snapshot_differences__test_table'
		SnapshotTables_suffix: 		'__snapshot__test_table'
		}
	Test_main()
		{
		cl = .snapshotTables
		cl.Drop(logTable?:)
		tableValues = .tableValues(tableNames = Object())

		// Test Ensure / test log table
		Assert(TableExists?(logTable = cl.SnapshotTables_logTable) is: false)
		cl.Ensure()
		Assert(TableExists?(logTable))

		// Attempt compare with no table snapshots
		differences = cl.Compare(tableNames)
		Assert(differences isSize: tableNames.Size())
		tableNames.Each()
			{
			Assert(differences has: 'No snapshot to compare against for table: ' $ it)
			}

		// Snapshot tables, then log the differences (none expected)
		cl(tableNames)
		differences = cl.Compare(tableNames)
		tableNames.Each()
			{
			snapshotTable = it $ cl.SnapshotTables_suffix
			Assert(TableExists?(snapshotTable))
			Assert(QueryCount(it) is: QueryCount(snapshotTable))
			}
		Assert(cl.Log(differences, persistSnaps?:) is: 0)
		Assert(QueryCount(logTable) is: 0)
		// Ensure persistSnaps?: argument prevents dropping the snapshot tables
		tableNames.Each({ Assert(TableExists?(it $ cl.SnapshotTables_suffix)) })
		// Ensure Drop drops all the snapshot tables
		Assert(cl.Drop() is: tableNames.Size())
		tableNames.Each({ Assert(TableExists?(it $ cl.SnapshotTables_suffix) is: false) })

		// Snapshot tables, modify their data, then log the differences (changes expected)
		cl(tableNames)
		tableValues.Each()
			{
			if it.Member?('changes')
				it.changes.Each({|change| QueryEnsure(it.table, change) })
			}

		// Ensure all snapshot tables are dropped when persistSnaps? is false
		Assert(cl.Log(cl.Compare(tableNames)) is: 9)
		Assert(review = cl.Review() isSize: 4)
		Assert(QueryAll(logTable) isSize: 9)
		tableNames.Each({ Assert(TableExists?(it $ cl.SnapshotTables_suffix) is: false) })

		// Assert that the logged differences are structured correctly
		// table1 should only have 1 tracked record, with one change
		Assert(table1Changes = review[tableValues.table1.table] isSize: 1)
		changeRec = table1Changes[#()]
		Assert(changeRec.diff.table1_data is: [snap: 'string', live: 'string '])

		// table2 should have 5 tracked records, with one change each
		Assert(table2Changes = review[tableValues.table2.table] isSize: 5)
		changeRec = table2Changes[#(table2_key: 0)]
		Assert(changeRec.diff.table2_data is: [snap: 'string', live: 'STRING'])
		changeRec = table2Changes[#(table2_key: 1)]
		Assert(changeRec.diff.table2_data is: [snap: false, live:])
		changeRec = table2Changes[#(table2_key: 2)]
		Assert(changeRec.diff.table2_data is: [snap: 10000, live: 10001])
		changeRec = table2Changes[#(table2_key: 3)]
		Assert(changeRec.diff.table2_data.snap isnt: changeRec.diff.table2_data.live)
		Assert(Date?(changeRec.diff.table2_data.snap))
		Assert(Date?(changeRec.diff.table2_data.live))
		changeRec = table2Changes[#(table2_key: 4)]
		Assert(changeRec.diff.table2_data is: [snap: [val: 1], live: [val: 2]])

		// table3 should have 1 tracked record with multiple changes
		table3Changes = review[tableValues.table3.table]
		Assert(table3Changes isSize: 1)
		changeRec = table3Changes[#(table3_key: 0)]
		Assert(changeRec.diff.table3_test? is: [snap:, live: false])
		Assert(changeRec.diff.table3_test!
			is: [snap: 'exclamation mark', live: 'exclamation'])

		// table4 should have 2 tracked records with multiple changes each
		table4Changes = review[tableValues.table4.table]
		Assert(table4Changes isSize: 2)
		changeRec = table4Changes[#(table4_key: 0)]
		Assert(changeRec.diff.table4_data is: [snap: 'test1', live: 'test1.1'])
		Assert(changeRec.diff hasMember: 'table4_TS')
		}

	tableValues(tableNames)
		{
		tables = Object(
			alwaysEmpty: [
				schema: '(empty_key, empty_data) key(empty_key)'
				],
			neverChanged: [
				schema: '(never_key, never_data) key(never_key)'
				data: [
					[never_key: 0, never_data: Timestamp()],
					[never_key: 1, never_data: Timestamp()]]],
			table1: [
				schema: '(table1_data) key()',
				data:    [[table1_data: 'string']],
				changes: [[table1_data: 'string ']]],
			table2: [
				schema: '(table2_key, table2_data) key(table2_key)',
				data: [
					[table2_key: 0, table2_data: 'string'],
					[table2_key: 1, table2_data: false],
					[table2_key: 2, table2_data: 10000],
					[table2_key: 3, table2_data: Date()],
					[table2_key: 4, table2_data: [val: 1]]],
				changes: [
					[table2_key: 0, table2_data: 'STRING'],
					[table2_key: 1, table2_data:],
					[table2_key: 2, table2_data: 10001],
					[table2_key: 3, table2_data: Date().Plus(seconds: 1)],
					[table2_key: 4, table2_data: [val: 2]]]],
			table3: [
				schema: '(table3_key, table3_data, table3_test?, table3_test!) ' $
					'key(table3_key) ' $
					'key(table3_data, table3_test!)',
				data:    [
					[table3_key: 0, table3_test?:, table3_test!: 'exclamation mark']]
				// Only one value changes between the records
				changes: [
					[table3_key: 0, table3_test?: false, table3_test!: 'exclamation']]],
			table4: [
				schema: '(table4_key, table4_data, table4_TS) key(table4_key)',
				data: [
					[table4_key: 0, table4_data: 'test1'],
					[table4_key: 1, table4_data: 'test2']]
				// Should also trigger the TS fields
				changes: [
					[table4_key: 0, table4_data: 'test1.1'],
					[table4_key: 1, table4_data: 'test2.2']]])

		tables.Each()
			{
			tableNames.Add(it.table = .MakeTable(it.schema))
			if it.Member?('data')
				it.data.Each({|rec| QueryOutput(it.table, rec) })
			}
		return tables
		}

	Test_Candidates()
		{
		cl = .snapshotTables
		table = .MakeTable('(table_key, table_data) key(table_key)')

		// Always excluded
		candidates = cl.Candidates()
		Object('tables', 'indexes', 'columns',
			'views', Test.TestLibName(), cl.SnapshotTables_logTable).
			Each({ Assert(candidates hasnt: it) })
		Assert(candidates has: table)
		}

	Test_formatTables()
		{
		m = SnapshotTables.SnapshotTables_formatTables
		Assert(m('') is: #(''))
		Assert(m('table1') is: #('table1'))
		Assert(m(#('table1')) is: #('table1'))
		Assert(m(#('table1', 'table2', 'table3')) is: #('table1', 'table2', 'table3'))

		all = m('all')
		Assert(all.Size() greaterThan: 0)
		Assert(all hasnt: 'tables')
		Assert(all hasnt: 'indexes')
		Assert(all hasnt: 'columns')
		Assert(all hasnt: 'views')
		Assert(all hasnt: .TestLibName())
		Assert(all hasnt: SnapshotTables.SnapshotTables_logTable)
		}

	Test_Pair()
		{
		cl = .snapshotTables
		cl.Drop(logTable?:)
		tableValues = .tableValues(tableNames = Object())

		// Test Ensure / test log table
		Assert(TableExists?(logTable = cl.SnapshotTables_logTable) is: false)
		cl.Ensure()
		Assert(TableExists?(logTable))

		// Snapshot tables, modify their data, then log the differences
		cl(tableNames)
		tableValues.Each()
			{
			if it.Member?('changes')
				it.changes.Each({|change| QueryEnsure(it.table, change) })
			}
		Assert(expectedPairs = cl.Log(cl.Compare(tableNames), persistSnaps?:) is: 9)

		// Ensure that we have the expected amount of live/snap pairs as per our snapshots
		pairs = 0
		QueryApply(logTable)
			{
			pair = cl.Pair(it.snap_table, it.snap_key_pairs)
			Assert(pair.snap isnt: pair.live)
			Assert(pair.snap isObject:)
			Assert(pair.live isObject:)
			pairs++
			}
		Assert(pairs is: expectedPairs)

		// Ensure we handle non-snapped tables
		Assert(cl.Pair('fake', #()) is: 'No snapshot to pair with table: fake')
		}

	Teardown()
		{
		.snapshotTables.Drop(logTable?:)
		super.Teardown()
		}
	}
