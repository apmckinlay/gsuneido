// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Setup()
		{
		.SvcTable(.LibExport = .MakeLibrary())
		.SvcLibrary = .SvcTable(.LibImport = .MakeLibrary())
		// Used by SVCLibIO_Test
		QueryOutput(.LibImport, [num: 999, group: 0, parent: 0, name: 'folder',
			lib_committed: #20210101])
		.SvcTable(.bookExport = .MakeBook())
		.SvcBook = .SvcTable(.BookImport = .MakeBook())
		}

	libIO_cl: LibIO
		{
		FileBlock(filename, block, mode)
			{
			Retry(maxRetries: 3, minDelayMs: 100)
				{
				super.FileBlock(:filename, :block, :mode)
				}
			}
		}

	VerifyExportImport(name, lib_committed, useSVCDates = false,
		lib_committed_expcted = '', modified_updated? = true)
		{
		_testStartDate = Date()
		.exportRecords(name, lib_committed, :useSVCDates)
		.assertRecords(name, '#test', lib_committed_expcted, :modified_updated?)
		}

	exportRecords(name, lib_committed, useSVCDates = false)
		{
		cl = new .libIO_cl()

		num = QueryMax(.LibExport, 'num', 0) + 1

		QueryOutput(.LibExport, Object(:num, :name,
			text: '#test', :lib_committed, group: -1, parent: 0))

		file = .MakeFile()
		cl.Export(.LibExport, name, file)
		QueryDo('delete ' $ .LibExport $ ' where num is ' $ num)
		cl.Import(file, .LibImport, :useSVCDates)

		QueryOutput(.bookExport, Object(num: 1, order: 10, :name,
			path: '', text: '#test', :lib_committed))
		file = .MakeFile()
		cl.Export(.bookExport, name, file, '')
		QueryDo('delete ' $ .bookExport $ ' where num is ' $ num)
		cl.Import(file, .BookImport, :useSVCDates)
		}

	assertRecords(name, text, lib_committed, modified_updated? = true)
		{
		rec = Query1(.LibImport, :name)
		Assert(rec.text.Trim() is: text.Trim())
		Assert(rec.lib_committed is: lib_committed)
		if modified_updated?
			Assert(rec.lib_modified greaterThanOrEqualTo: _testStartDate)
		else
			Assert(rec.lib_modified is: '')

		rec = Query1(.BookImport, :name, path: '')
		Assert(rec.text.Trim() is: text.Trim())
		Assert(rec.lib_committed is: lib_committed)
		if modified_updated?
			Assert(rec.lib_modified greaterThanOrEqualTo: _testStartDate)
		else
			Assert(rec.lib_modified is: '')
		}

	VerifyExportImportForDelete(name, lib_committed, hasDeleted = false,
		useSVCDates = false)
		{
		.exportDelete(name, lib_committed, hasDeleted, useSVCDates)
		.assertDelete(name, lib_committed, useSVCDates)
		}

	exportDelete(name, lib_committed, hasDeleted, useSVCDates)
		{
		cl = new .libIO_cl()

		rec = [:name, text: "#text", :lib_committed, path: 'path']
		if not hasDeleted
			{
			.SvcLibrary.Output(rec.Copy().Merge([parent: 0]))
			.SvcBook.Output(rec.Copy().Merge([order: 10]))
			Assert(Query1(.LibImport, :name, group: -1) isnt: false)
			Assert(Query1(.BookImport, :name, path: 'path') isnt: false)
			}
		else
			{
			.SvcLibrary.Output(rec.Copy(), deleted:)
			.SvcBook.Output(rec.Copy(), deleted:)
			Assert(Query1(.LibImport, :name, group: -2) isnt: false)
			Assert(Query1(.BookImport, :name, path: '<deleted>path') isnt: false)
			}

		file = .MakeFile()
		Assert({ cl.Export(.LibExport, name, file) }
			throws: "LibIO Export can't get " $ name $ " from " $ .LibExport)
		cl.Export(.LibExport, name, file, delete:)
		cl.Import(file, .LibImport, :useSVCDates)

		file = .MakeFile()
		Assert({ cl.Export(.bookExport, name, file, 'path') }
			throws: "LibIO Export can't get " $ name $ " from " $ .bookExport)
		cl.Export(.bookExport, name, file, 'path', delete:)
		cl.Import(file, .BookImport, :useSVCDates)
		}

	assertDelete(name, lib_committed, useSVCDates)
		{
		Assert(Query1(.LibImport, :name, group: -1) is: false)
		Assert(Query1(.BookImport, :name, path: 'path') is: false)
		if useSVCDates
			{
			Assert(Query1(.LibImport, :name) is: false)
			Assert(Query1(.BookImport, :name, path: '<deleted>path') is: false)
			}
		else
			{
			deleteRec = Query1(.LibImport, :name)
			Assert(deleteRec.lib_committed is: lib_committed)
			deleteRec = Query1(.BookImport, :name, path: '<deleted>path')
			Assert(deleteRec.lib_committed is: lib_committed)
			}
		}
	}
