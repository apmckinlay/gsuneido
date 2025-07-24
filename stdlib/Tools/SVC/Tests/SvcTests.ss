// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	LibraryCacheOverride: 	#SvcDisabledLibrariesOverride
	BookCacheOverride: 		#SvcDisabledBooksOverride
	MakeLibrary(@records)
		{
		library = super.MakeLibrary(@records)
		.TearDownIfTablesNotExist(library $ '_master')
		.tableCacheOverride(library, .LibraryCacheOverride)
		return library
		}

	tableCacheOverride(table, cache)
		{
		if Suneido.Member?(cache)
			Suneido[cache].Add(table)
		else
			Suneido[cache] = Object(table)
		}

	MakeBook()
		{
		book = super.MakeBook()
		.TearDownIfTablesNotExist(book $ '_master')
		.tableCacheOverride(book, .BookCacheOverride)
		return book
		}

	MakeMasterTable(type, table = false)
		{
		if table is false
			table = .TempTableName()
		SvcCore.EnsureMaster(masterTable = table $ '_master', type)
		.AddTeardown({ Database('drop ' $ masterTable) })
		return masterTable
		}

	svcTable: SvcTable
		{
		ResetSvcDisabledCache()
			{
			cache = .Type is 'book'
				? SvcTests.BookCacheOverride
				: SvcTests.LibraryCacheOverride
			Suneido.GetDefault(cache, Object()).Remove(super.Table())
			super.ResetSvcDisabledCache()
			}
		}
	SvcTable(table)
		{
		return .svcTable(table, svcEnsure:)
		}

	svc: Svc
		{
		ProcessFeedbackOb(unused) { }
		Svc_deleteRecordContrib(@unused) {}
		}
	Svc(server = false, local? = true)
		{
		svc = .svc(server, local?)
		if svc.Svc_svc.Base() is SvcClient
			throw 'ERROR: Svc: Tests should NOT run with SvcClient'
		return svc
		}

	svcModel: SvcModel
		{
		New(.testParent)
			{ }
		SvcModel_initSvc()
			{
			settings = .SvcModel_settings
			return .testParent.Svc(settings.server, settings.local?)
			}
		}
	SvcModel(table)
		{
		model = new .svcModel(this)
		model.SetSettings([server: false, local?:])
		model.SetTable(table)
		return model
		}

	AssertSvcEmpty(svc, table)
		{
		Assert(svc.Master_changes(table) is: #())
		Assert(svc.Local_changes(table) is: #())
		Assert(svc.Conflicts(table) is: #())
		}

	CommitAdd(svc, svcTable, name, text, desc, id = #default, lib_before_hash = '',
		path = '')
		{
		table = svcTable.Table()
		QueryOutput(table, [
			:name, :text, :lib_before_hash, :path,
			num: NextTableNum(table),
			group: -1,
			lib_modified: Timestamp()
			])
		if svcTable.Type is 'book'
			name = path $ '/' $ name
		Assert(svc.Put(svcTable, name, id, desc) isDate: true)
		return svc.Get(table, name).lib_committed
		}

	CommitTextChange(svc, svcTable, name, text, desc, id = #default)
		{
		table = svcTable.Table()
		QueryDo('update ' $ table $ ' where name is ' $ Display(name) $ ' set text = ' $
			Display(text) $ ', lib_modified = ' $ Display(Timestamp()))
		Assert(svc.Put(svcTable, name, id, desc) isDate: true)
		return svc.Get(table, name).lib_committed
		}

	CommitDelete(svc, svcTable, name, desc, id = #default)
		{
		svcTable.StageDelete(name)
		svc.Remove(svcTable, name, id, desc)
		return svc.GetDel(svcTable.Table(), name).lib_committed
		}

	OutstandingDelete(svc, svcTable, name, text, desc, id = #default)
		{
		ts = .CommitAdd(svc, svcTable, name, text, 'Added: ' $ desc)
		svc.Remove(svcTable, name, id, 'Deleted: ' $ desc)
		if svcTable.GetMaxCommitted() > ts
			svcTable.SetMaxCommitted(ts, force:)
		}

	OutstandingConflict(svc, svcTable, name, text, lib_committed = false)
		{
		table = svcTable.Table()
		if lib_committed is false
			lib_committed = .getDateFromPrevRecord(svc, table, name)
		QueryApply1(table, :name, group: -1)
			{
			it.text = text
			it.lib_modified = Timestamp()
			it.lib_committed = lib_committed
			it.Update()
			}
		if lib_committed is "" or svcTable.GetMaxCommitted() > lib_committed
			svcTable.SetMaxCommitted(lib_committed, force:)
		}

	getDateFromPrevRecord(svc, table, name)
		{
		// Get the current master record, at this point, it maybe identical to the local
		if false is rec = svc.Get(table, name)
			throw 'ERROR: cannot find master record to conflict with, (' $ name $ ')'
		// Get the previous record
		prevRec = svc.GetBefore(table, name, rec.lib_committed)
		// prevRec should be either a different record OR false
		Assert(prevRec isnt: rec)
		return prevRec isnt false ? prevRec.lib_committed : ''
		}

	Teardown()
		{
		for cache in [.LibraryCacheOverride, .BookCacheOverride]
			if Suneido.Member?(cache)
				Suneido.Delete(cache)
		super.Teardown()
		}
	}
