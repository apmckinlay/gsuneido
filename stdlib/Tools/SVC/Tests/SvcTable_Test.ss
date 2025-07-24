// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_Get()
		{
		svcTableBk = .SvcTable(.MakeBook())
		svcTableBk.Output([name: #one, parent: 0, order: 1])
		svcTableBk.Output([name: #two, parent: 0, order: 2])

		svcTableLib1 = .SvcTable(tbl1 = .MakeLibrary())
		svcTableLib1.Output([name: #one, parent: 0])
		svcTableLib1.Output([name: #two, parent: 0])

		svcTableLib2 = .SvcTable(tbl2 = .MakeLibrary())
		svcTableLib2.Output([name: #one, parent: 0])

		// setting it so NextTableNum is higher
		QueryDo('update ' $ tbl2 $ ' where name is #one set num = 10')
		svcTableLib2.Output([name: #two, parent: 0])

		Assert(svcTableLib2.Get(#one).num is: 10)
		Assert(svcTableLib2.Get(#two).num is: 11)

		Assert(svcTableLib1.Get(#one, table: tbl1).num is: 1)
		Assert(svcTableLib1.Get(#two, table: tbl1).num is: 2)

		Assert(svcTableBk.Get(#one).num is: 1)
		Assert(svcTableBk.Get(#two).num is: 2)
		}

	Test_DeleteRestore()
		{
		svc = .setup('DeleteRestore')

		// Records do not exist, ensuring code coverage
		// SvcLibrary
		.svcTableLb.StageDelete(.libRecName)
		// SvcBook
		.svcTableBk.StageDelete(.bookRecName)

		// Record exists, not committed
		// SvcLibrary
		.svcTableLb.Output(.libRec.Copy())
		.assertDelete(.svcTableLb, .libRecName)
		// SvcBook
		.svcTableBk.Output(.bookRec.Copy())
		.assertDelete(.svcTableBk, .bookRecName)

		// Records exist and are committed to SVC
		.svcTableLb.Output(.libRec.Copy())
		.svcTableBk.Output(.bookRec.Copy())
		// Add Svc required fields to Library record
		libRec = .libRec.Copy()
		libRec.lib = .lib
		libRec.type = '+'
		// Add Svc required fields to Book record
		bookRec = .bookRec.Copy()
		bookRec.lib = .book
		bookRec.type = '+'
		bookRec.name = .bookRecName
		changes = [libRec, bookRec]
		svc.SendLocalChanges(changes, 'committing records', 'SvcTable_Test')
		// Delete the record, confirm it is in the delete table
		// SvcLibrary
		.assertDelete(.svcTableLb, .libRecName, stagedDelete:)
		Assert(del = .svcTableLb.Get(.libRecName, deleted:) isnt: false)
		Assert(del.lib_before_text is: 'DeleteRestore')
		Assert(del.lib_before_path is: '/Fake')
		// SvcBook
		.assertDelete(.svcTableBk, .bookRecName, stagedDelete:)
		Assert(del = .svcTableBk.Get(.bookRecName, deleted:) isnt: false)
		Assert(del.lib_before_text is: 'Order: 7\r\n\r\nDeleteRestore')
		Assert(del.lib_before_path is: '/Fake')

		// Restore the deleted records
		// SvcLibrary
		Assert(.svcTableLb.Get(.libRecName) is: false)
		.svcTableLb.Restore(.libRecName)
		Assert(.svcTableLb.Get(.libRecName, deleted:) is: false)
		Assert(restored = .svcTableLb.Get(.libRecName) isnt: false)
		Assert(restored.text is: 'DeleteRestore')
		Assert(restored.lib_before_text is: '')
		Assert(restored.lib_before_path is: '')
		// SvcBook
		Assert(.svcTableBk.Get(.bookRecName) is: false)
		.svcTableBk.Restore(.bookRecName)
		Assert(.svcTableBk.Get(.bookRecName, deleted:) is: false)
		Assert(restored = .svcTableBk.Get(.bookRecName) isnt: false)
		Assert(restored.order is: 7)
		Assert(restored.text is: 'Order: 7\r\n\r\nDeleteRestore')
		Assert(restored.lib_before_text is: '')
		Assert(restored.path is: '/Fake')
		Assert(restored.lib_before_path is: '')
		}

	setup(purpose, parent = 'Fake', name = false)
		{
		svc = .Svc()

		// Create test tables / variables
		// Library
		.lib = .MakeLibrary()
		// "()" are to test that Library Use state does not cause issues
		.svcTableLb = .SvcTable('(' $ .lib $ ')')
		name = name is false ? .TempName() : name
		.MakeLibraryRecord([num: num = NextTableNum(.lib), group: 0, name: parent],
			table: .lib)
		.libRec = [:name, group: -1, parent: num, text: purpose, path: '/' $ parent]
		.libRecName = .svcTableLb.MakeName(.libRec)

		// Book
		.book = .MakeBook()
		.svcTableBk = .SvcTable(.book)
		.bookRec = .libRec.Copy()
		.bookRec.order = 7
		.bookRecName = .svcTableBk.MakeName(.bookRec)

		return svc
		}

	assertDelete(svcTable, recName, stagedDelete = false)
		{
		Assert(svcTable.Get(recName) isnt: false)
		svcTable.StageDelete(recName)
		Assert(svcTable.Get(recName) is: false)
		if stagedDelete
			Assert(svcTable.Get(recName, deleted:) isnt: false)
		else
			Assert(svcTable.Get(recName, deleted:) is: false)
		}

	Test_RestoreModified()
		{
		.setup('RestoreModified')
		.MakeLibraryRecord([num: NextTableNum(.lib), group: 0, name: 'Orig'],
			table: .lib)

		// Records do not exist, testing no issues occur
		// SvcLibrary
		.svcTableLb.Restore(.libRecName)
		// SvcBook
		.svcTableBk.Restore(.bookRecName)

		// Record exists, no lib_committed, records are deleted
		// SvcLibrary
		.svcTableLb.Output(.libRec.Copy())
		.svcTableLb.Restore(.libRecName)
		Assert(.svcTableLb.Get(.libRecName) is: false)
		Assert(.svcTableLb.Get(.libRecName, deleted:) is: false)
		// SvcBook
		.svcTableBk.Output(.bookRec.Copy())
		.svcTableBk.Restore(.bookRecName)
		Assert(.svcTableBk.Get(.bookRecName) is: false)
		Assert(.svcTableBk.Get(.bookRecName, deleted:) is: false)

		// Records exists, has lib_committed, lib_modified is '', no changes occur
		// SvcLibrary
		.libRec.lib_committed = Date()
		.svcTableLb.Output(.libRec.Copy())
		QueryDo('update ' $ .lib $ ' where name is "' $ .libRec.name $ '"
			set lib_modified = ""')
		rec = .svcTableLb.Get(.libRecName)
		.svcTableLb.Restore(.libRecName)
		Assert(.svcTableLb.Get(.libRecName) is: rec)
		// Book
		.bookRec.lib_committed = Date()
		.svcTableBk.Output(.bookRec.Copy())
		QueryDo('update ' $ .book $ ' where name is "' $ .bookRec.name $ '"
			set lib_modified = ""')
		rec = .svcTableBk.Get(.bookRecName)
		.svcTableBk.Restore(.bookRecName)
		Assert(.svcTableBk.Get(.bookRecName) is: rec)

		// Records exists, has lib_committed, lib_modified isnt '', restore occurs
		// SvcLibrary
		QueryApply1(.svcTableLb.NameQuery(.libRecName))
			{
			.svcTableLb.Update(it, it.Transaction(), newText: 'Modified Text')
			}
		.svcTableLb.Restore(.libRecName)
		rec = .svcTableLb.Get(.libRecName)
		Assert(rec.text is: 'RestoreModified')
		// Book
		.bookRec.lib_modified = Date()
		QueryApply1(.svcTableBk.NameQuery(.bookRecName))
			{
			.svcTableBk.Update(it, it.Transaction(), newText: 'Modified Text')
			}
		.svcTableBk.Restore(.bookRecName)
		// Page name is changed due to Restore updating .path
		rec = .svcTableBk.Get('/Fake/' $ .bookRec.name)
		Assert(rec.text is: 'Order: 7\r\n\r\nRestoreModified')
		Assert(rec.order is: 7)
		}

	Test_Rename()
		{
		.setup('Rename', 'Parent', 'Rename')
		.svcTableBk.Output(rec = [name: #Parent, parent: 0, text: 'Parent Text'])

		// Record exists, rename occurs, new records, original records are not stored
		// SvcLibrary
		.svcTableLb.Output(.libRec.Copy())
		.assertRename(.svcTableLb, 'Rename', 'Rename2')
		// SvcBook
		.svcTableBk.Output(rec = .bookRec.Copy())
		.svcTableBk.Output(
			[name: #Child1, parent: rec.num, text: 'Child1 Text', path: '/Parent/Rename']
			)
		.svcTableBk.Output(
			[name: #Child2, parent: rec.num, text: 'Child2 Text',
				path: '/Parent/Rename/Child1']
			)
		.assertRename(.svcTableBk, '/Parent/Rename', 'Rename2')

		// Record exists, rename occurs, committed records, original record is staged for
		// deletion
		// SvcLibrary
		QueryDo('update ' $ .svcTableLb.NameQuery('Rename2') $
			' set lib_committed = ' $ Display(Date()))
		.assertRename(.svcTableLb, 'Rename2', 'Rename3', stagedDelete:)
		// SvcBook
		QueryDo('update ' $ .svcTableBk.NameQuery('/Parent/Rename2') $
			' set lib_committed = ' $ Display(Date()))
		.assertRename(.svcTableBk, '/Parent/Rename2', 'Rename3', stagedDelete:)

		// Record is renamed to deleted record
		// SvcLibrary
		.assertRename(.svcTableLb, 'Rename3', 'Rename2', restoreDelete:)
		// SvcBook
		.assertRename(.svcTableBk, '/Parent/Rename3', 'Rename2', restoreDelete:)
		}

	assertRename(svcTable, oldName, newName, stagedDelete = false, restoreDelete = false)
		{
		Assert(renameRec = svcTable.Get(oldName) isnt: false)
		svcTable.Rename(renameRec.Copy(), newName)
		Assert(svcTable.Get(oldName) is: false)
		if stagedDelete
			Assert(svcTable.Get(oldName, deleted:) isnt: false)
		else
			Assert(svcTable.Get(oldName, deleted:) is: false)
		// Exclude record name from comparison as they will naturally differ
		newRec = svcTable.Get(svcTable.MakeName([name: newName, path: renameRec.path]))
		if restoreDelete
			{
			Assert(newRec.text is: renameRec.text)
			Assert(newRec.path is: renameRec.path)
			// Ensure lib_committed is populate (from the deleted record)
			Assert(newRec.lib_committed isDate: true)
			newRec.Delete(#lib_before_text, #lib_before_path, #lib_committed)
			}
		// Ensure renamed book records update children records
		if svcTable.Type is 'book'
			{
			Assert(svcTable.Get(svcTable.MakeName(newRec) $ '/Child1') isnt: false)
			Assert(svcTable.Get(svcTable.MakeName(newRec) $ '/Child1/Child2') isnt: false)
			}
		}

	Test_MoveLibrary()
		{
		lib = .MakeLibrary()
		svcTable = .SvcTable(lib)

		svcTable.Output([name: 'Location1', group: 0])
		loc1 = Query1(lib, name: 'Location1')

		svcTable.Output([name: 'Location2', group: 0])
		loc2 = Query1(lib, name: 'Location2')

		svcTable.Output([name: 'Location3', group: loc2.num, parent: loc2.num,
			path: 'Location2'])
		loc3 = Query1(lib, name: 'Location3')

		name = .TempName()
		libRec = [:name, parent: loc1.num, text: 'Move']
		Assert(libRec.path = svcTable.GetPath(libRec) is: 'Location1')
		svcTable.Output(libRec)

		svcTable.Move(Query1(lib, num: num = libRec.num), loc3.num)
		Assert(rec = svcTable.Get(name) isnt: false)
		Assert(rec.path is: 'Location2/Location3')
		Assert(rec.lib_before_path is: 'Location1')
		Assert(rec.lib_modified isnt: '')

		svcTable.Move(Query1(lib, :num), loc2.num)
		Assert(rec = svcTable.Get(name) isnt: false)
		Assert(rec.path is: 'Location2')
		Assert(rec.lib_before_path is: 'Location1')
		Assert(rec.lib_modified isnt: '')

		svcTable.Move(Query1(lib, :num), loc1.num)
		Assert(rec = svcTable.Get(name) isnt: false)
		Assert(rec.path is: 'Location1')
		Assert(rec.lib_before_path is: '')
		}

	Test_MoveBook()
		{
		book = .MakeBook()
		svcTable = .SvcTable(book)
		name = .TempName()
		svcTable.Output(loc1 = [name: 'Location1', path: ''])
		svcTable.Output(loc2 = [name: 'Location2', path: ''])
		svcTable.Output(loc3 = [name: 'Location3', path: '/' $ loc2.name])

		bookRec = [:name, path: '/' $ loc1.name, text: 'Move']
		svcTable.Output(bookRec)
		num = bookRec.num

		// Record is not yet committed. Deleted record is NOT staged for svc deletion
		// during move, it is simply deleted
		svcTable.Move(Query1(book, :num), loc3.num)
		newName = svcTable.MakeName([path: svcTable.MakeName(loc3), name: bookRec.name])
		Assert(rec = svcTable.Get(newName) isnt: false)
		Assert(rec.path is: '/' $ loc2.name $ '/' $ loc3.name)
		prevName = svcTable.MakeName([path: rec.path, name: bookRec.name])
		Assert(svcTable.Get(prevName, deleted:) is: false)

		// Record is now committed. Deleted record IS staged for svc deletion during move
		QueryDo('update ' $ book $ ' where num is ' $ rec.num $
			' set lib_committed = ' $ Display(Date()))
		svcTable.Move(Query1(book, :num), loc2.num)
		newName = svcTable.MakeName([path: svcTable.MakeName(loc2), name: bookRec.name])
		Assert(rec = svcTable.Get(newName) isnt: false)
		Assert(rec.path is: '/' $ loc2.name)
		prevName = svcTable.MakeName([path: rec.path, name: bookRec.name])
		Assert(svcTable.Get(prevName, deleted:) is: false)

		// Now that the record is new again, it is NOT staged for svc deletion,
		// it is simply deleted
		svcTable.Move(Query1(book, :num), loc1.num)
		newName = svcTable.MakeName([path: svcTable.MakeName(loc1), name: bookRec.name])
		Assert(rec = svcTable.Get(newName) isnt: false)
		Assert(rec.path is: '/' $ loc1.name)
		prevName = svcTable.MakeName([path: rec.path, name: bookRec.name])
		Assert(svcTable.Get(prevName, deleted:) is: false)

		// Record is moved back to committed position. Pull info from deleted rec.
		// Remove deleted record
		svcTable.Move(Query1(book, :num), loc3.num)
		newName = svcTable.MakeName([path: svcTable.MakeName(loc3), name: bookRec.name])
		Assert(rec = svcTable.Get(newName) isnt: false)
		Assert(rec.path is: '/' $ loc2.name $ '/' $ loc3.name)
		prevName = svcTable.MakeName([path: rec.path, name: bookRec.name])
		Assert(svcTable.Get(prevName, deleted:) is: false)
		}

	Test_MaxCommitted()
		{
		for table in [.MakeLibrary(), .MakeBook()]
			{
			svcTable = .SvcTable(table)

			// Empty table, max commit record is output with: #17000101
			.assertMaxCommitted(svcTable, Date.Begin())

			// Delete record repeat process, this time with a record in the table.
			// Max commit record should be output with the tables max date
			QueryDo('delete ' $ table)
			svcTable.Output([name: 'Rec0', lib_committed: Date()], committed:)
			svcTable.Output(
				[name: 'Rec1', lib_committed: date = Date().Plus(hours: 1)], committed:)
			.assertMaxCommitted(svcTable, date)

			// Delete record test set process, with a record in the table.
			// Ensure that SetMaxCommitted outputs the max commit record using the tables
			// max lib_committed value (as the date specified is less then)
			QueryDo('delete ' $ table $ ' where name is "' $ svcTable.MaxCommitName $ '"')
			svcTable.SetMaxCommitted(Date())
			.assertMaxCommitted(svcTable, date, recExists:)

			// Delete record test set process, with a record in the table.
			// Ensure that SetMaxCommitted outputs the max commit record with the
			// date specified
			QueryDo('delete ' $ table $ ' where name is "' $ svcTable.MaxCommitName $ '"')
			svcTable.SetMaxCommitted(date = Date().Plus(hours: 2))
			.assertMaxCommitted(svcTable, date, recExists:)

			// Ensure that SetMaxCommitted updates the max commit record with the
			// date specified
			svcTable.SetMaxCommitted(date = Date().Plus(hours: 3))
			.assertMaxCommitted(svcTable, date, recExists:)

			// Ensure that SetMaxCommitted DOES NOT update the max commit record with the
			// date specified as it is less then the current value
			svcTable.SetMaxCommitted(Date())
			.assertMaxCommitted(svcTable, date, recExists:)

			// Delete all records from table, ensure that the set outputs the record
			// and uses the date specified
			QueryDo('delete ' $ table)
			svcTable.SetMaxCommitted(date = Date())
			.assertMaxCommitted(svcTable, date, recExists:)

			// Esnure "force" parameter forces the Set to set the date
			Assert(svcTable.GetMaxCommitted() is: date)
			svcTable.SetMaxCommitted(Date.Begin())
			Assert(svcTable.GetMaxCommitted() is: date)
			svcTable.SetMaxCommitted(Date.Begin(), force:)
			Assert(svcTable.GetMaxCommitted() is: Date.Begin())
			}
		}

	assertMaxCommitted(svcTable, date, recExists = false)
		{
		if not recExists
			Assert(Query1(svcTable.Table(), name: svcTable.MaxCommitName) is: false)
		Assert(svcTable.GetMaxCommitted() is: date)
		Assert(rec = Query1(svcTable.Table(), name: svcTable.MaxCommitName) isnt: false)
		if svcTable.Type is #book
			Assert(rec.path is: #commitRec)
		else
			{
			Assert(rec.group is: -3)
			Assert(rec.parent is: 0)
			}
		}

	Test_DeleteHandling()
		{
		svcTable = .SvcTable(.MakeLibrary())
		rec1 = .deletedRec('Rec1')
		svcTable.Output(rec1, deleted:)
		Assert(result = QueryAll(svcTable.Query(deleted:)) isSize: 1)
		Assert(result has: rec1)

		rec2 = .deletedRec('Rec2')
		svcTable.Output(rec2, deleted:)
		Assert(result = QueryAll(svcTable.Query(deleted:)) isSize: 2)
		Assert(result has: rec1)
		Assert(result has: rec2)

		rec3 = .deletedRec('Rec3')
		svcTable.Output(rec3, deleted:)
		Assert(result = QueryAll(svcTable.Query(deleted:)) isSize: 3)
		Assert(result has: rec1)
		Assert(result has: rec2)
		Assert(result has: rec3)
		}

	deletedRec(name, lib_before_text = '', lib_before_path = '', lib_committed = '')
		{
		return [:name, :lib_committed, :lib_before_text, :lib_before_path].
			RemoveIf({ it is '' })
		}
	}