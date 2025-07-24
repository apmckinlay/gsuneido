// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_update_move()
		{
		bookModel = BookEditModel(book = .MakeBook())
		svcTable = .SvcTable(book)
		folder1 =.MakeBookRecord(book, 'Folder1',
			extrafields: #(lib_committed: #17000101.0010))
		folder2 = .MakeBookRecord(book, 'Folder2',
			extrafields: #(lib_committed: #17000101.0020))

		rec = .MakeBookRecord(book, 'Item1', path: '/' $ folder1.name,
			extrafields: #(lib_committed: #17000101.0011))

		bookModel.Move(bookModel.Get(rec.num), folder2.num)
		newRec = Query1(book, num: rec.num)
		Assert(newRec.lib_committed is: '')
		Assert(newRec.lib_modified isDate: true)
		Assert(newRec.lib_before_path is: '/' $ folder1.name)
		Assert(newRec.path is: '/' $ folder2.name)
		Assert(QueryCount(svcTable.Query(deleted:))  is: 1)

		bookModel.Move(bookModel.Get(rec.num), folder1.num)
		newRec = Query1(book, name: rec.name)
		Assert(newRec.path is: '/' $ folder1.name)
		Assert(QueryCount(svcTable.Query(deleted:)) is: 0)
		Assert(newRec.lib_committed is: #17000101.0011)
		}

	Test_update_rename()
		{
		bookModel = BookEditModel(book = .MakeBook())
		svcTable = .SvcTable(book)
		item = .MakeBookRecord(book, 'Item',
			extrafields: #(lib_committed: #17000101.0011))

		// Rename record
		oldName = item.name
		bookModel.Rename(bookModel.Get(item.num), newName = .TempName())

		// Delete record is output and a new record is created
		del = svcTable.Get('/' $ oldName, deleted:)
		Assert(del.lib_committed is: #17000101.0011)
		Assert(del.lib_before_text is: 'Order: \r\n\r\nItem')
		Assert(del.lib_before_path is: '')

		Assert(item = bookModel.Get(item.num) isnt: false)
		Assert(item.name is: newName)
		Assert(item.lib_committed is: '')
		Assert(item.text is: 'Item')

		// Rename record back to original name
		bookModel.Rename(bookModel.Get(item.num), oldName)
		// No deleted records are output, original record is no longer staged for svc
		// deletion
		Assert(QueryCount(svcTable.Query(deleted:)) is: 0)
		Assert(item = bookModel.Get(item.num) isnt: false)
		Assert(item.name is: oldName)
		Assert(item.lib_committed is: #17000101.0011)
		Assert(item.lib_before_text is: 'Order: \r\n\r\nItem')
		Assert(item.lib_before_path is: '')
		}

	Test_Synced?()
		{
		m = BookEditModel.Synced?

		rec = [lib_committed: '', lib_modified: '']
		savedRec = [lib_committed: '', lib_modified: '']
		Assert(m(rec, savedRec))

		rec.lib_committed = Date()
		Assert(m(rec, savedRec) is: false)

		savedRec.lib_committed = rec.lib_committed
		Assert(m(rec, savedRec))

		rec.lib_modified = Date()
		Assert(m(rec, savedRec) is: false)

		savedRec.lib_modified = rec.lib_modified
		Assert(m(rec, savedRec))

		rec.lib_committed = savedRec.lib_committed
		rec.lib_modified = ''
		Assert(m(rec, savedRec) is: false)

		rec.lib_modified = savedRec.lib_modified
		Assert(m(rec, savedRec))

		rec.text = '<p>test text</p>'
		Assert(m(rec, savedRec) is: false)

		savedRec.text = '<p>test text</p>'
		Assert(m(rec, savedRec))

		savedRec.text = '<p>Test text</p>'
		Assert(m(rec, savedRec) is: false)
		}

	Test_Update()
		{
		query1CacheLogs = .SpyOn(Query1CacheReset).Return('').CallLogs()
		.SvcTable(book = .MakeBook())
		becl = BookEditModel { ClearCache () { } }
		model = becl(book)

		// Updating record which doesn't exist, ensure no issues arise
		// Name is HtmlPrefix, so Query1CacheReset is run
		rec = [name: 'HtmlPrefix']
		model.Update(rec)
		Assert(query1CacheLogs isSize: 1)

		// Updating record which does exist
		// Name is HtmlSuffix, so Query1CacheReset is run
		rec = .MakeBookRecord(book, 'start text 1',
			path: .TempName() $ '/' $ .TempName())
		rec.name = 'HtmlSuffix'
		QueryDo('update ' $ book $ ' where num is ' $ rec.num $
			' set name = "HtmlSuffix"')
		Assert(Query1(book, num: rec.num).text is: 'start text 1')
		model.Update([num: rec.num, text: 'end text 1', name: 'HtmlSuffix', order: 10])
		Assert(rec = Query1(book, num: rec.num) isnt: false)
		Assert(rec.text is: 'end text 1')
		Assert(rec.order is: 10)
		Assert(query1CacheLogs isSize: 2)

		// Updating record which does exist
		// Name isnt HtmlPrefix or HtmlSuffix so Query1CacheReset isnt run
		rec = .MakeBookRecord(book, 'start text 2',
			path: .TempName() $ '/' $ .TempName())
		Assert(Query1(book, num: rec.num).text is: 'start text 2')
		model.Update([num: rec.num, text: 'end text 2'])
		Assert(rec = Query1(book, num: rec.num) isnt: false)
		Assert(rec.text is: 'end text 2')
		Assert(rec.order is: '')
		Assert(query1CacheLogs isSize: 2)

		// Making record seem "unmodified" to test modifying record Order
		QueryApply1(book, num: rec.num)
			{
			it.lib_committed = Date()
			it.lib_modified = it.lib_before_text = it.lib_before_path = ''
			it.order = 1
			it.Update()
			rec = it
			}
		// Record has no changes from what is saved, lib_modified stays ''
		model.Update(rec)
		Assert(Query1(book, num: rec.num).lib_modified is: '')
		// Order is the only thing that has changed, lib_modified is set
		rec.order = 7
		model.Update(rec)
		rec = Query1(book, num: rec.num)
		Assert(rec.lib_modified isDate:)
		Assert(rec.order is: 7)
		}

	Test_TreeSort()
		{
		rec1 = [name: 'B', order: 1]
		rec2 = [name: 'A', order: 2]
		mock = Mock(BookEditModel)
		mock.When.Get(rec1).Return(rec1)
		mock.When.Get(rec2).Return(rec2)
		mock.When.TreeSort([anyArgs:]).CallThrough()

		// rec1 comes first based on order
		Assert(mock.TreeSort(rec1, rec2) is: -1)

		// rec1 still comes first
		rec1.name = 'AA'
		Assert(mock.TreeSort(rec1, rec2) is: -1)

		// rec2 comes first as both recs share the same order, (falls back to alpahetical)
		rec2.order = 1
		Assert(mock.TreeSort(rec1, rec2) is: 1)
		}

	Test_copyName()
		{
		m = BookEditModel.BookEditModel_copyName

		Assert(m('Test Page') is: s = 'Test Page Copy 1')
		Assert(m(s) is: s = 'Test Page Copy 2')
		Assert(m(s) is: s = 'Test Page Copy 3')
		Assert(m(s) is: s = 'Test Page Copy 4')
		Assert(m(s) is: s = 'Test Page Copy 5')

		Assert(m('Ends With Copy') is: s = 'Ends With Copy Copy 1')
		Assert(m(s) is: s = 'Ends With Copy Copy 2')
		Assert(m(s) is: s = 'Ends With Copy Copy 3')
		Assert(m(s) is: s = 'Ends With Copy Copy 4')
		Assert(m(s) is: s = 'Ends With Copy Copy 5')

		Assert(m('Double Digits Copy 9') is: s = 'Double Digits Copy 10')
		Assert(m(s) is: s = 'Double Digits Copy 11')
		Assert(m(s) is: s = 'Double Digits Copy 12')

		Assert(m('Triple Digits Copy 99') is: s = 'Triple Digits Copy 100')
		Assert(m(s) is: s = 'Triple Digits Copy 101')
		Assert(m(s) is: s = 'Triple Digits Copy 102')
		}
	}
