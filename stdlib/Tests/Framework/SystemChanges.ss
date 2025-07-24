// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.

// used by TestRunner to detect when a test creates/removes tables or files
class
	{
	GetState()
		{
		return [tables: .TableStates(), files: Dir(details:)]
		}

	// returns "" if no differences
	CompareState(before, after = false)
		{
		if after is false
			after = .GetState()
		return .compareState(before, after)
		}

	compareState(before, after)
		{
		return Object(
			.checkTableDifferences(before.tables, after.tables),
			.checkFileDifferences(before.files, after.files)).
				Remove('').Join('\n')
		}

	// tables ------------------------------------------------------------------

	TableStates(_systemChanges_excludeTables = false)
		{
		if not Object?(systemChanges_excludeTables)
			systemChanges_excludeTables = Object()
		systemChanges_excludeTables.AddUnique('Test_lib')

		if Sys.Client?()
			return ServerEval('SystemChanges.TableStates', systemChanges_excludeTables)

		tableStates = Object()
		// Distancing from the next Timestamp to better detect changes to _TS values
		asof = Date().Replace(millisecond: 0).Minus(seconds: 1)
		QueryApply('tables')
			{
			tableStates[it.table] = [nrows: it.nrows, totalsize: it.totalsize, :asof]
			}
		systemChanges_excludeTables.Each({ .excludeTable(tableStates, it) })
		return tableStates
		}

	excludeTable(tableStates, table)
		{
		if not tableStates.Member?(table)
			return
		where = ' where table is ' $ Display(table)
		tableStates.indexes.nrows -= QueryCount('indexes' $ where)
		tableStates.columns.nrows -= QueryCount('columns' $ where)
		tableStates.tables.nrows--
		tableStates.Delete(table)
		}

	checkTableDifferences(tables_before, tables_after)
		{
		tableDifferencesOb = Object()
		tables_before.Members().MergeUnion(tables_after.Members()).Each()
			{
			// If false, then the table is new
			beforeStats = tables_before.GetDefault(it, false)
			// If false, then the table was deleted
			afterStats = tables_after.GetDefault(it, false)
			if '' isnt msg = .tableDifference(it, beforeStats, afterStats)
				tableDifferencesOb.Add(it $ ': ' $ msg)
			}
		return Opt('Table Discrepancies:\r\n\t- ', tableDifferencesOb.Join('\r\n\t- '))
		}

	tableDifference(table, beforeStats, afterStats)
		{
		if beforeStats is false or afterStats is false
			return beforeStats is false
				? 'created'
				: 'deleted'

		differences = Object()
		if 0 isnt nrowsDiff = afterStats.nrows - beforeStats.nrows
			differences.Add('nrows: ' $ nrowsDiff)

		if 0 isnt sizeDiff = .checkSize(table, beforeStats, afterStats)
			differences.Add('totalsize: ' $ sizeDiff)
		return differences.Join(', ')
		}

	checkSize(table, beforeStats, afterStats)
		{
		return table not in ('tables', 'columns', 'indexes')
			? afterStats.totalsize - beforeStats.totalsize
			: 0
		}

	// files -------------------------------------------------------------------

	checkFileDifferences(files_before, files_after)
		{
		excludedFiles = GetContributions('TestExcludedFiles')
		// handle when the test copy is running inside tmp folder
		tmpFolder = GetAppTempPath.SubFolderName()
		excludedFiles.Add(tmpFolder $ '/')
		equal = .removeExcludedFiles(files_before, excludedFiles) is
					.removeExcludedFiles(files_after, excludedFiles)
		return equal ? "" : .buildErrorString(files_before, files_after)
		}

	removeExcludedFiles(files, excludedFiles)
		{
		return files.RemoveIf(
			{ |f| excludedFiles.Any?({ |exclude| f.name.Prefix?(exclude)}) })
		}

	buildErrorString(files_before, files_after)
		{
		return .newFiles(files_after, files_before) $
				.changedFiles(files_before, files_after)
		}

	changedFiles(files_before, files_after)
		{
		errorStr = ""
		diffOb = files_before.Difference(files_after)
		for file in diffOb
			{
			beforeMem = .findFile(files_before, file.name)
			afterMem = .findFile(files_after, file.name)
			if beforeMem is false
				continue
			if afterMem is false
				{
				errorStr $= 'DELETED ' $ file.name
				continue
				}
			errorStr $= "CHANGED " $ file.name $ ": "
			fileBefore = files_before[beforeMem]
			fileAfter = files_after[afterMem]
			for mem in file.Members()
				if fileBefore[mem] isnt fileAfter[mem]
					errorStr $= mem $ ": before " $
						Display(fileBefore[mem]) $ ", after " $
						Display(fileAfter[mem]) $ "; "
			}
		return errorStr
		}

	newFiles(files_after, files_before)
		{
		errorStr = ""
		newOb = files_after.Difference(files_before)
		for file in newOb
			if false is .findFile(files_before, file.name)
				errorStr $= "NEW " $ file.name $ ", size: " $ file[#size] $ "; "
		return errorStr
		}

	findFile(files, name)
		{
		return files.FindIf({ it.name is name })
		}

	}
