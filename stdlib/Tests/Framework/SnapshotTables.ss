// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
/* USAGE
	SnapshotTables will copy the data from one table into a separate/basic copy of it.
	The copy will only use one key, and will not maintain any of the table's relations.
	Additionally, all of the columns will be renamed to avoid triggering rules.

	These "snaps" can be used to compare or track modifications made to a table.
	For example:
		table = 'example_table'
		SnapshotTables.Snap(table) 		// Creates a snapshot of: example_table
		Example_Test()					// Test is ran, making changes to: example_table
		SnapshotTables.Compare(table) 	// The snapshot and live table are compared

	You can either review the returned object directly, or you can output it to a table.
	To do this, do the following:
		table = 'example_table'
		SnapshotTables.Ensure() 					// Ensure the log table is ready/empty
		differences = SnapshotTables.Compare(table)	// Collect the differences
		SnapshotTables.Log(differences)				// Log the differences accordingly

	Normally, the snapped tables will be dropped once they are logged.
	If you want to keep the snapped tables after logging them,
	simply pass in: persistSnaps?: true

	If you have records in snapshot_differences, and the "snaps" are still available,
	you can use .Pair(<snap_table>, <snap_key_pairs>) to retrieve the two specific records
	for review.
*/
class
	{
	CallClass(tables = '(All)')
		{
		.formatTables(tables).Each(.Snap)
		}

	formatTables(tables)
		{
		if Object?(tables)
			return tables
		return tables is '(All)'
			? .Candidates()
			: Object(tables)
		}

	suffix: '__snapshot'
	Candidates()
		{
		excludeTables = Object('tables', 'indexes', 'columns',
			'views', Test.TestLibName(), .logTable).Map!(Display).Join(', ')
		return QueryList('tables
			where
				table not in (' $ excludeTables $ ') and
				not table.Suffix?(`' $ .suffix $'`)',
			'table').Sort!()
		}

	Snap(table)
		{
		snapshot = .ensureSnap(table)
		QueryApplyMulti(.query(table), update:)
			{
			QueryOutput(snapshot, it)
			}
		}

	ensureSnap(table)
		{
		try Database('drop ' $ snapshot = table $ .suffix)
		columns = '(' $ QueryColumns(table).Map(.rename).Join(', ') $ ') '
		keys = 'key (' $ .keys(table).Map(.rename).Join(', ') $ ')'
		Database('ensure ' $ snapshot $ columns $ keys)
		return snapshot
		}

	keys(table)
		{
		return ShortestKey(table).Split(',').Map!({ it.Trim() }).Remove('')
		}

	questionMarkSub: `__qu_mark`
	exclamationMarkSub: `__ex_mark`
	rename(column)
		{
		if column.Suffix?(`?`)
			column = column.RemoveSuffix(`?`) $ .questionMarkSub
		else if column.Suffix?(`!`)
			column = column.RemoveSuffix(`!`) $ .exclamationMarkSub
		return column $ .suffix
		}

	query(table)
		{
		renames = Object()
		QueryColumns(table).Each({ renames.Add(it $ ' to ' $ .rename(it)) })
		return table $ Opt(' rename ', renames.Join(', '))
		}

	Compare(tables = '(All)')
		{
		differences = Object()
		.formatTables(tables).Each({ differences[it] = .compare(it) })
		return differences
		}

	compare(table)
		{
		if not TableExists?(snapshotTable = table $ .suffix)
			return 'No snapshot to compare against for table: ' $ table

		columns = .potentialColumns(table, snapshotTable)
		differences = Object().Set_default(Object())
		.iterateSnapshotPair(table, snapshotTable)
			{|snapshotRec, liveRec, where|
			columns.Each()
				{|rename|
				snap = snapshotRec[rename]
				live = liveRec[rename]
				if snap isnt live
					differences[where][.reverseRename(rename)] = [:live, :snap]
				}
			}
		return differences
		}

	potentialColumns(table, snapshotTable)
		{
		return QueryColumns(table).Map(.rename).MergeUnion(QueryColumns(snapshotTable))
		}

	iterateSnapshotPair(table, snapshotTable, block)
		{
		keys = .keys(snapshotTable)
		QueryApply(snapshotTable)
			{|snapshotRec|
			keyPairs = Object()
			where = Object()
			keys.Each()
				{
				keyPairs[.reverseRename(it)] = snapshotRec[it]
				where.Add(it $ ' is ' $ Display(snapshotRec[it]))
				}
			QueryApply1(.query(table) $ Opt(' where ', where.Join(' and ')))
				{
				block(snapshotRec, it, keyPairs)
				}
			}
		}

	reverseRename(rename)
		{
		original = rename.RemoveSuffix(.suffix)
		if original.Suffix?(.questionMarkSub)
			return original.RemoveSuffix(.questionMarkSub) $ `?`
		if original.Suffix?(.exclamationMarkSub)
			return original.RemoveSuffix(.exclamationMarkSub) $ `!`
		return original
		}

	Drop(logTable? = false)
		{
		return .drop(.Snaps(), logTable?)
		}

	Snaps()
		{
		return QueryList('tables where table.Suffix?(`' $ .suffix $ '`)', 'table')
		}

	drop(tables, logTable?)
		{
		if logTable?
			tables.Add(.logTable)
		dropped = 0
		tables.Each()
			{
			try
				{
				Database('drop ' $ it)
				dropped++
				}
			}
		return dropped
		}

	Log(tableDifferences, persistSnaps? = false)
		{
		output = 0
		for snap_table, differences in tableDifferences
			if Object?(differences)
				for snap_key_pairs, snap_differences in differences
					{
					QueryOutput(.logTable,
						[:snap_table, :snap_key_pairs, :snap_differences])
					output++
					}
		if not persistSnaps?
			.drop(tableDifferences.Members().Map!({ it $ .suffix }), false)
		return output
		}

	Review()
		{
		review = Object().Set_default(Object())
		QueryApply(.logTable)
			{
			review[it.snap_table][it.snap_key_pairs] =
				Object(
					diff: it.snap_differences,
					pair: .Pair(it.snap_table, it.snap_key_pairs))
			}
		return review
		}

	logTable: 'snapshot_differences'
	Ensure()
		{
		try Database('drop ' $ .logTable)
		Database('ensure ' $ .logTable $
			' (snap_table, snap_key_pairs, snap_differences) ' $
			' key (snap_table, snap_key_pairs)')
		}

	Pair(table, snap_key_pairs)
		{
		if not TableExists?(snapshotTable = table $ .suffix)
			return 'No snapshot to pair with table: ' $ table
		wheres = .keyPairWheres(snap_key_pairs)
		snap = Query1(snapshotTable $ Opt(' where ', wheres.snap))
		live = Query1(table $ Opt(' where ', wheres.live))
		return [:snap, :live]
		}

	keyPairWheres(snap_key_pairs)
		{
		snap = Object()
		live = Object()
		for key, value in snap_key_pairs
			{
			isValue = ' is ' $ Display(value)
			snap.Add(key $ .suffix $ isValue)
			live.Add(key $ isValue)
			}
		return [snap: snap.Join(' and '), live: live.Join(' and ')]
		}
	}
