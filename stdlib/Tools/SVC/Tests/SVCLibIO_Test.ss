// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
LibIOTestBase
	{
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
	svcLibIO_cl: SVCLibIO
		{
		FileBlock(filename, block, mode)
			{
			Retry(maxRetries: 3, minDelayMs: 100)
				{
				super.FileBlock(:filename, :block, :mode)
				}
			}
		}
	fileContents: '{nameHolder}
librec_info: #(orig_lib: "{libHolder}", lib_committed: {commitHolder}, delete:)

==========================================================='
	Test_Main()
		{
		.VerifyExportImportForDelete([name: "rec1", group: -1, parent: 999, type: '-',
			lib_committed: #20210520], useSVCDates: true, expectedCommit: #20210520)
		// Not from SVC Export, don't update lib_committed
		.VerifyExportImportForDelete([name: "rec2", group: -1, parent: 999, type: '-',
			lib_committed: #20210521], useSVCDates: false, expectedCommit: #20210520)
		// update lib_committed
		.VerifyExportImportForDelete([name: "rec3", group: -1, parent: 999, type: '-',
			lib_committed: #20210528], useSVCDates: true, expectedCommit: #20210528)
		}

	VerifyExportImportForDelete(record, useSVCDates = false, expectedCommit = false)
		{
		.exportDelete(record, useSVCDates)
		.assertDelete(record, useSVCDates, expectedCommit)
		}

	exportDelete(record, useSVCDates)
		{
		cl = new .svcLibIO_cl()

		record.num = QueryMax(.LibImport, 'num', 0) + 1
		QueryOutput(.LibImport, record)
		Assert(Query1(.LibImport, name: record.name, group: -1) isnt: false)


		file = .MakeFile(.setNameLibCommit(record, .fileContents))

		cl.Export(.LibExport, record, file, delete:)
		cl.Import(file, .LibImport, quiet:, :useSVCDates)
		}

	setNameLibCommit(record, fileString)
		{
		return  fileString.Replace("{libHolder}", .LibExport).
			Replace("{nameHolder}", record.name).
			Replace('{commitHolder}', Display(record.lib_committed))
		}

	assertDelete(record, useSVCDates, expectedCommit)
		{
		Assert(Query1(.LibImport, name: record.name, group: -1) is: false)
		res = Query1(.LibImport, name: record.name, group: -2)
		if useSVCDates
			Assert(res is: false)
		else
			Assert(res isnt: false)
		if expectedCommit isnt false
			Assert(SvcTable(.LibImport).GetMaxCommitted() is: expectedCommit)
		}

	Test_export_import_book_master()
		{
		.SvcTable(book = .MakeBook())
		masterTable = .MakeMasterTable('book')
		QueryOutput(masterTable,[name: '/Payables/Apply Checks', path: '/Payables',
			text: 'Order: 1\r\n\r\nAp_ApplyPaymentControl',
			lib_committed: #20080101])
		QueryOutput(masterTable, [name: '/Payables/Help Doc 1', path: '/Payables',
			text: 'Order: 2\r\n\r\nA help text document',
			lib_committed: #20080101])
		QueryOutput(masterTable, [name: '/Payables/Help Doc 2', path: '/Payables',
			text: 'Order: \r\n\r\nA help text document with no order',
			lib_committed: #20080101])
		QueryOutput(masterTable, [name: '/Intorduction', path: '',
			text: 'Order: 1\r\n\r\nPath is at the root',
			lib_committed: #20080101])
		QueryOutput(masterTable, [name: '/End', path: '',
			text: 'Order: 10\r\n\r\n',
			lib_committed: #20080101])

		_testStartDate = Date()
		file = .MakeFile()
		.export(masterTable, '/Apply Checks', '/Payables', file)
		.export(masterTable, '/Help Doc 1', '/Payables', file)
		.export(masterTable, '/Help Doc 2', '/Payables', file)
		.export(masterTable, '/Intorduction', '', file)
		.export(masterTable, '/End', '', file)
		// Cannot use SvcLibIO for importing records.
		// SvcLibIO.Get_record expects a record to be passed to it in place of "name".
		// During the export, this occurs, but during the import the real "name" is used
		Assert(.libIO_cl.Import(file, book, useSVCDates:) is: 5)

		.assertBookRecord(book, "Apply Checks", '/Payables', 'Ap_ApplyPaymentControl',
			#20080101, order: 1)
		.assertBookRecord(book, "Help Doc 1", '/Payables', 'A help text document',
			#20080101, order: 2)
		.assertBookRecord(book, "Help Doc 2", '/Payables',
			'A help text document with no order', #20080101)
		.assertBookRecord(book, "Intorduction", '', 'Path is at the root',
			#20080101, order: 1)
		.assertBookRecord(book, "End", '', '', #20080101, order: 10)
		}

	export(masterTable, name, path, file)
		{
		rec = Query1(masterTable, name: path $ name, :path)
		rec.name = name.AfterLast('/')
		SvcBook.SplitText(rec)
		.svcLibIO_cl.Export(masterTable, rec, file, :path)
		}

	assertBookRecord(book, name, path, text, lib_committed, order = '')
		{
		rec = Query1(book, :name, :path)
		Assert(rec.text is: text)
		Assert(rec.lib_committed is: lib_committed)
		Assert(rec.order is: order)
		Assert(rec.lib_modified is: '')
		}
	}
