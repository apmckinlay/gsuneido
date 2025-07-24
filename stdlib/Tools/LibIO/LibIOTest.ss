// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
SvcTests
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

	Test_main()
		{
		.table1 = .MakeTable('(name, text, num, group, parent,
			lib_committed, lib_modified)
			key(name)')
		.table2 = .MakeTable('(name, text, num, group, parent,
			lib_committed, lib_modified)
			key(name)')

		.outputRecords()
		file = .MakeFile()
		.libIO_cl.Export(.table2, 'Init', file)
		.libIO_cl.Export(.table2, 'Window', file)

		// testing lib_modified is "" from import
		file2 = .MakeFile()
		.libIO_cl.Export(.table2, 'Test_lib_keepSVCDates', file2)

		// blank lines at end of file should have no effect on the Import
		AddFile(file, '\r\n\r\n\r\n\r\n\r\n\r\n')

		// two records should be imported
		Assert(.libIO_cl.Import(file, .table1) is: 2)
		fileRec = Query1(.table1, name: "Window")
		Assert(Date?(fileRec.lib_modified))
		Assert(fileRec.lib_committed is: "")

		x = Query1(.table2, name: "Window")
		Assert(x isnt: false)
		y = Query1(.table1, name: "Window")
		Assert(y isnt: false)
		Assert(y.num is: 4)
		// should be in "TestFolder2" folder which has num = 2 in .table1
		Assert(y.parent is: 2)
		Assert(x.name is y.name and x.text is y.text and y.lib_committed is "")

		x = Query1(.table2, name: 'Init')
		Assert(x isnt: false)
		y = Query1(.table1, name: 'Init')
		Assert(y isnt: false)
		Assert(y.num is: 3)
		Assert(y.parent is: 0)
		Assert(x.name is y.name and x.text is y.text and
			x.lib_committed is y.lib_committed)

		Assert(QueryDo("update " $ .table1 $ " where group is -1 set text = 123") is: 2)

		// two records should be imported
		Assert(.libIO_cl.Import(file, .table1) is: 2)
		y = Query1(.table1, name: 'Init')
		Assert(Date?(y.lib_modified))
		Assert(y isnt: false, msg: "y isnt False")
		Assert(x.name is y.name and x.text is y.text, "name and text preserved")

		// one record imported with lib_modified empty
		Assert(.libIO_cl.Import(file2, .table1, useSVCDates:) is: 1)
		fileRec2 = Query1(.table1, name: "Test_lib_keepSVCDates")
		Assert(fileRec2.lib_modified is: '')
		Assert(fileRec2.lib_committed is: #20080101)
		}

	outputRecords()
		{
init_text =
"function ()
	{
	Testing()
	}
"
window_text =
"class
	{
	New()
		{
		}
	}
"

		QueryOutput(.table1, Object(num: 1, parent: 0, group: 0,
			name: 'TestFolder', text: ''))
		QueryOutput(.table1, Object(num: 2, parent: 0, group: 0,
			name: 'TestFolder2', text: ''))

		QueryOutput(.table2, Object(num: 1, parent: 0, group: 0,
			 name: 'TestFolder2', text: ''))
		QueryOutput(.table2, Object(num: 2, parent: 1, group: -1,
			name: "Window", text: window_text,
			lib_committed: #20080101))
		QueryOutput(.table2, Object(num: 3, parent: 0, group: -1,
			name: "Init", text: init_text))
		QueryOutput(.table2, Object(num: 4, parent: 0, group: -1,
			name: "Test_lib_keepSVCDates", text: init_text
			lib_committed: #20080101))
		}

	Test_ImportIntoMultiLibs()
		{
		.table4 = .MakeTable('(name, text, num, group, parent, lib_invalid_text,
			lib_committed, lib_modified)
			key(name)')
		.table5 = .MakeTable('(name, text, num, group, parent,
			lib_committed, lib_modified)
			key(name)')
		tableNonExist = .MakeTable('(name, text, num, group, parent,
			lib_committed, lib_modified)
			key(name)')
		QueryOutput(.table4, Object(num: 1, parent: 0, group: -1,
			name: "File1", text: 'Text 1'))
		QueryOutput(.table5, Object(num: 1, parent: 0, group: -1,
			name: "File2", text: 'Text 2'))
		QueryOutput(tableNonExist, Object(num: 1, parent: 0, group: -1,
			name: "File3", text: 'Text 3'))
		QueryOutput(.table4, Object(num: 1, parent: 0, group: -1,
			name: "File5", text: 'Valid Code', lib_invalid_text: 'Invalid Code'))

		file = .MakeFile()
		.libIO_cl.Export(.table4, 'File1', file)
		.libIO_cl.Export(.table5, 'File2', file)
		.libIO_cl.Export(tableNonExist, 'File3', file)
		Assert({ .libIO_cl.Export(tableNonExist, 'File4', file) }
			throws: 'LibIO Export can\'t get File4 from '$ tableNonExist)
		.libIO_cl.Export(tableNonExist, 'File4', file, delete:)
		Assert({ .libIO_cl.Export(.table4, 'File5', file) }
			throws: 'LibIO did not export: ' $ .table4 $ ':File5' $
				' as it contains errors. Please correct and re-attempt export')

		Database('destroy ' $ tableNonExist)

		QueryDelete(.table4, Object(name: 'File1', group: -1))
		QueryDelete(.table5, Object(name: 'File2', group: -1))
		Assert(Query1(.table4, name: 'File1', group: -1) is: false)
		Assert(Query1(.table5, name: 'File2', group: -1) is: false)

		// Test import with lib = false
		// Records should be imported into the libs they were exported from
		.libIO_cl.Import(file, false)
		Assert(Query1(.table4, name: 'File1', group: -1).text.Trim() is: 'Text 1')
		Assert(Query1(.table5, name: 'File2', group: -1).text.Trim() is: 'Text 2')

		QueryDelete(.table4, Object(name: 'File1', group: -1))
		QueryDelete(.table5, Object(name: 'File2', group: -1))
		Assert(Query1(.table4, name: 'File1', group: -1) is: false)
		Assert(Query1(.table5, name: 'File2', group: -1) is: false)

		// Test import with lib specified.
		// Records should be imported into the specified lib
		.libIO_cl.Import(file, .table4)
		Assert(Query1(.table4, name: 'File1', group: -1).text.Trim() is: 'Text 1')
		Assert(Query1(.table4, name: 'File2', group: -1).text.Trim() is: 'Text 2')
		Assert(Query1(.table5, name: 'File2', group: -1) is: false)
		}

	Test_unload()
		{
		lib = .MakeLibraryRecord(Record(name: "LibIOTestValue", text: "1"))
		file = .MakeFile()
		.libIO_cl.Export(lib, 'LibIOTestValue', file)
		QueryDo("update " $ lib $ " set text = '2'")
		Unload('LibIOTestValue')
		// use Global to avoid code warnings since normally not defined
		Assert(Global('LibIOTestValue') is: 2)
		.libIO_cl.Import(file, lib)
		Assert(Global('LibIOTestValue') is: 1)
		}

	Test_path()
		{
		.table3 = .MakeTable('(name, text, num, group, parent,
			lib_committed, lib_modified)
			key(name)')
		QueryOutput(.table3, Object(num: 1, parent: 0, group: 0,
			name: 'Level1', text: ''))
		QueryOutput(.table3, Object(num: 2, parent: 1, group: 1,
			name: 'Level2', text: ''))

		QueryOutput(.table3, Object(num: 3, parent: 2, group: -1,
			name: "TestRec1", text: 'function() {/* do nothing */}'))

		file = .MakeFile()
		.libIO_cl.Export(.table3, 'TestRec1', file)

		str = GetFile(file)
		ob = str.AfterFirst('librec_info: ').BeforeFirst('\r\n').SafeEval()
		Assert(ob hasMember: 'path', msg: 'Path missing from export')

		QueryDelete(.table3, Record(name: 'TestRec1', num: 3))

		// Test Basic Import
		Assert(QueryEmpty?(.table3, num: 3), msg: 'table3 not empty')
		.libIO_cl.Import(file, .table3)
		rec = QueryFirst(.table3 $ ' where num is 3 sort name')
		Assert(rec isnt: false)
		Assert(rec.parent is: 2)

		// Test Missing Level
		QueryDelete(.table3, Record(name: 'TestRec1', num: 3))
		QueryDelete(.table3, Record(name: 'Level2', num: 2))

		.libIO_cl.Import(file, .table3)
		parentRec = Query1(.table3, name: 'Level2', group: 1)
		Assert(parentRec isnt: false)
		Assert(Query1(.table3, name: 'TestRec1', parent: parentRec.num) isnt: false)

		// Test whole path missing
		QueryDelete(.table3, Record(name: 'TestRec1', num: 3))
		QueryDelete(.table3, Record(name: 'Level2', num: 2))
		QueryDelete(.table3, Record(name: 'Level1', num: 1))

		.libIO_cl.Import(file, .table3)
		parentRec = Query1(.table3, name: 'Level1', group: 0)
		Assert(parentRec isnt: false)
		parentRec = Query1(.table3, name: 'Level2', group: parentRec.num)
		Assert(parentRec isnt: false)
		Assert(Query1(.table3, name: 'TestRec1', parent: parentRec.num) isnt: false)
		}

	Test_add_book_recinfo()
		{
		lib_record = Record()
		Assert(.libIO_cl.LibIO_add_book_recinfo(lib_record) is: '')

		lib_record.lib_modified = #20080101.1200
		Assert(.libIO_cl.LibIO_add_book_recinfo(lib_record)
			is: '// lib_modified: #20080101.1200, lib_committed: ')

		lib_record.lib_committed = #20080101.1201
		Assert(.libIO_cl.LibIO_add_book_recinfo(lib_record)
			is: '// lib_modified: #20080101.1200, lib_committed: #20080101.1201')

		lib_record.lib_modified = ''
		Assert(.libIO_cl.LibIO_add_book_recinfo(lib_record)
			is: '// lib_modified: , lib_committed: #20080101.1201')

		lib_record.text = '<h2>Test</h2>'
		Assert(.libIO_cl.LibIO_add_book_recinfo(lib_record)
			is: '<!-- lib_modified: , lib_committed: #20080101.1201-->')

		lib_record.lib_modified = #20080102.1333
		Assert(.libIO_cl.LibIO_add_book_recinfo(lib_record)
			is: '<!-- lib_modified: #20080102.1333, lib_committed: #20080101.1201-->')
		}

	Test_extractBookRecInfo()
		{
		book_info = Object(text: '')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info is: #(text: ''))

		book_info = Object(text: 'test')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info is: #(text: 'test'))

		book_info = Object(text: 'Taxes
// lib_modified: , lib_committed: #20020626.081452940
')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info.text is: 'Taxes')
		Assert(book_info.lib_modified is: "")
		Assert(book_info.lib_committed is: #20020626.081452940)

		book_info = Object(text: 'Taxes
// lib_modified: #20020626.0900, lib_committed: #20020626.081452940
')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info.text is: 'Taxes')
		Assert(book_info.lib_modified is: #20020626.0900)
		Assert(book_info.lib_committed is: #20020626.081452940)

		book_info = Object(text: 'Taxes
// lib_modified: #20020626.0900, lib_committed:
')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info.text is: 'Taxes')
		Assert(book_info.lib_modified is: #20020626.0900)
		Assert(book_info.lib_committed is: '')

		book_info = Object(text: '<h2>Test</h2>
<!-- lib_modified: #20020626.0900, lib_committed: -->')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info.text is: '<h2>Test</h2>')
		Assert(book_info.lib_modified is: #20020626.0900, msg: 'HTML modified 1')
		Assert(book_info.lib_committed is: '', msg: 'HTML committed 1')

		book_info = Object(text: '<h2>Test</h2>
<!-- lib_modified: #20020626.0900, lib_committed: #20020627.1333-->')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info.text is: '<h2>Test</h2>')
		Assert(book_info.lib_modified is: #20020626.0900, msg: 'HTML modified 2')
		Assert(book_info.lib_committed is: #20020627.1333, msg: 'HTML committed 2')

		book_info = Object(text: '
// lib_modified: #20131102.1130, lib_committed: #20131101.1130')
		.libIO_cl.LibIO_extractBookRecInfo(book_info)
		Assert(book_info.text is: '')
		Assert(book_info.lib_modified is: #20131102.1130)
		Assert(book_info.lib_committed is: #20131101.1130)
		}

	Test_export_import_book()
		{
		.SvcTable(book = .MakeBook())
		QueryOutput(book, Object(num: 1, order: 10, name: "Apply Checks",
			path: '/Payables', text: 'Ap_ApplyPaymentControl',
			lib_committed: #20080101, lib_modified: #20080101.1100))
		QueryOutput(book, Object(num: 2, name: "res",
			path: '', text: '', lib_committed: #20080102))
		QueryOutput(book, Object(num: 3, name: "empty_folder",
			lib_before_text: 'old text',
			path: '/res', text: '', lib_committed: #20080103))
		QueryOutput(book, Object(num: 4, name: "item",
			path: '/res/empty_folder', text: '//testing',
			lib_committed: #20080104, lib_modified: #20080104.1300))
		QueryOutput(book, Object(num: 5, name: "item2",
			path: '/res/empty_folder', text: '//testing\r\nitem2',
			lib_modified: #20080105))

		_testStartDate = Date()
		file = .MakeFile()
		.libIO_cl.Export(book, 'Apply Checks', file, path: '/Payables')
		.libIO_cl.Export(book, 'res', file, path: '')
		.libIO_cl.Export(book, 'empty_folder', file, path: '/res')
		.libIO_cl.Export(book, 'item', file, path: '/res/empty_folder')
		Assert(.libIO_cl.Import(file, book) is: 4)

		.assertBookRecord(book, "Apply Checks", '/Payables',
			'Ap_ApplyPaymentControl', #20080101, order: 10)
		.assertBookRecord(book, "res", '', '', #20080102)
		.assertBookRecord(book, "empty_folder", '/res', '', #20080103,
			lib_before_text: 'old text')
		.assertBookRecord(book, "item", '/res/empty_folder', '//testing', #20080104)

		// item2 should not be affected
		rec = Query1(book, name: "item2", path: '/res/empty_folder')
		Assert(rec.text is: '//testing\r\nitem2')
		Assert(rec.lib_committed is: '')
		Assert(rec.lib_modified is: #20080105)
		}

	assertBookRecord(book, name, path, text, lib_committed, modified_updated? = true,
		order = '', lib_before_text = false)
		{
		rec = Query1(book, :name, :path)
		Assert(rec.text is: text)
		Assert(rec.lib_before_text is:
			lib_before_text is false
				? 'Order: ' $ order $ '\r\n\r\n' $ text
				: lib_before_text)
		Assert(rec.lib_committed is: lib_committed)
		Assert(rec.order is: order)

		if modified_updated?
			Assert(rec.lib_modified greaterThanOrEqualTo: _testStartDate)
		else
			Assert(rec.lib_modified is: '')
		}

	Test_export_import_book_res()
		{
		text = ""
		for i in .. 256
			text $= i.Chr()
		Assert(text isSize: 256)
		cksum = Adler32(text)
		book = .MakeBook()
		QueryOutput(book, Object(num: 1, order: 10, name: "test.gif",
			path: '/res', :text,
			lib_committed: #20100805))
		s = Query1(book, name: 'test.gif').text
		Assert(s is: text msg: 'comparing output text')
		Assert(Adler32(s) is: cksum msg: 'comparing output text checksum')
		file = .MakeFile()
		.libIO_cl.Export(book, 'test.gif', file, path: '/res')
		Assert(Adler32(Query1(book, name: 'test.gif').text) is: cksum)
		.libIO_cl.Import(file, book)
		t = Query1(book, name: 'test.gif').text
		Assert(t isSize: text.Size(), msg: 'comparing imported text size')
		Assert(t is: text msg: 'comparing imported text')
		Assert(Adler32(t) is: cksum msg: 'comparing imported text checksum')
		}

	Test_import?()
		{
		mock = Mock(LibIO)
		mock.When.import?([anyArgs:]).CallThrough()
		mock.When.skipImport?([anyArgs:]).CallThrough()
		mock.When.skipImport([anyArgs:]).Do({ })

		// Can't find record, record is imported
		mock.When.Get_record([anyArgs:]).Return(false)
		imported = Object()
		result = mock.LibIO_import?('Testlib', 'TestRec', [text: 'this is a test'],
			imported)
		Assert(result)
		Assert(imported hasMember: 'Testlib:TestRec')
		Assert(imported['Testlib:TestRec'])
		mock.Verify.Get_record([anyArgs:])

		// Record is already marked for import, don't check Get_record again
		result = mock.LibIO_import?('Testlib', 'TestRec', [text: 'this is a test'],
			imported)
		Assert(result)
		mock.Verify.Get_record([anyArgs:])

		// Record is found, text matches import file, record is skipped
		mock.When.Get_record([anyArgs:]).Return([text: 'this is a test'])
		result = mock.LibIO_import?('Testlib', 'TestRec1', [text: ' this is a test '],
			imported)
		Assert(result is false)
		Assert(imported hasntMember: 'Testlib:TestRec1')
		mock.Verify.skipImport([anyArgs:])

		// Record is found, text matches import file, order is different, record is not skipped
		mock.When.Get_record([anyArgs:]).Return([text: 'another test', order: 1.0])
		result = mock.LibIO_import?('Testlib', 'TestRec1a',
			[text: 'another test', order: 1.1], imported)
		Assert(result)
		Assert(imported hasMember: 'Testlib:TestRec1a')

		// Record is not marked to be skipped, text differs, import occurs
		result = mock.LibIO_import?('Testlib', 'TestRec1', [text: ''], imported)
		Assert(result)
		mock.Verify.skipImport([anyArgs:])
		Assert(imported hasMember: 'Testlib:TestRec1')
		Assert(imported['Testlib:TestRec1'])

		// Record is found, text differs, record is not modified, import occurs
		mock.When.Get_record([anyArgs:]).Return([text: 'this is a test'])
		result = mock.LibIO_import?('Testlib', 'TestRec2', [text: ''], imported)
		Assert(result)
		Assert(imported hasMember: 'Testlib:TestRec2')
		Assert(imported['Testlib:TestRec2'])

		// Record is found, text differs, record is modified, askOverwrite returns false
		// import not carried out
		mock.When.askOverwrite([anyArgs:]).Return(false)
		mock.When.Get_record([anyArgs:]).
			Return([text: 'this is a test', lib_modified: Date()])
		result = mock.LibIO_import?('Testlib', 'TestRec3', [text: ''], imported)
		Assert(result is: false)
		Assert(imported hasMember: 'Testlib:TestRec3')
		Assert(imported['Testlib:TestRec3'] is: false)
		mock.Verify.askOverwrite([anyArgs:])

		// Record is already marked to be skipped, don't askOverwrite again
		result = mock.LibIO_import?('Testlib', 'TestRec3', [text: ''], imported)
		Assert(result is: false)
		mock.Verify.askOverwrite([anyArgs:])

		// Record is found, text differs, record is modified, askOverwrite returns true
		// import carried out
		mock.When.askOverwrite([anyArgs:]).Return(true)
		mock.When.Get_record([anyArgs:]).
			Return([text: 'this is a test', lib_modified: Date()])
		result = mock.LibIO_import?('Testlib', 'TestRec4', [text: ''], imported)
		Assert(result)
		Assert(imported hasMember: 'Testlib:TestRec4')
		Assert(imported['Testlib:TestRec4'])
		mock.Verify.Times(2).askOverwrite([anyArgs:])

		// Record is already marked to be imported, don't askOverwrite again
		result = mock.LibIO_import?('Testlib', 'TestRec4', [text: ''], imported)
		Assert(result)
		mock.Verify.Times(2).askOverwrite([anyArgs:])
		}
	}
