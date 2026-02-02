// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_changesToSend()
		{
		extra_lib = .MakeLibrary()
		excluded_library = .MakeLibrary([name: 'excluded_test', lib_modified: Date()])
		.SvcTable(book = .MakeBook())
		.SvcTable(lib = .MakeLibrary())
		Suneido.LibraryTablesOverride = Object(lib, excluded_library, extra_lib)
		Suneido.BookTablesOverride = Object(book)
		.SpyOn(SvcControl.Getter_SvcExcludeLibraries).Return(Object(excluded_library))

		func = SvcControl.SvcControl_changesToSend
		Assert(func() is: #())

		.MakeLibraryRecord([num: 1, name: 'test_rec', text: '"testing"',
			lib_modified: "", lib_committed: #20010101], table: lib)
		Assert(func() is: #()) // empty lib_modified

		.MakeLibraryRecord([num: 2, name: 'test_rec2', text: '"testing"',
			lib_modified: #20010102, lib_committed: #20010101], table: lib)

		.MakeBookRecord(book, 'testing book',
			extrafields: #(lib_modified: "", lib_committed: #20010101))
		Assert(func() is: Object(lib)) // no lib_modified on book record

		.MakeBookRecord(book, 'testing book2',
			extrafields: #(lib_modified: #20110102, lib_committed: #20010101))
		result = func()
		Assert(result isSize: 2)
		Assert(result has: lib)
		Assert(result has: book)

		.MakeLibraryRecord([name: 'deleted_rec', lib_committed: #19000101, group: -2],
			table: extra_lib)
		result = func()
		Assert(result isSize: 3)
		Assert(result has: lib)
		Assert(result has: book)
		Assert(result has: extra_lib)
		}

	Test_refreshRequired()
		{
		mock = Mock(SvcControl)
		mock.When.refreshRequired([anyArgs:]).CallThrough()

		// Cannot find local record and not a delete, refresh list to remove it
		mock.SvcControl_model = FakeObject(GetLocalRec: false)
		Assert(mock.refreshRequired([]))
		Assert(mock.refreshRequired([svc_type: '+']))
		Assert(mock.refreshRequired([svc_type: ' ']))
		Assert(mock.refreshRequired([svc_type: '-']))

		// Local pre-existing record, has no changes, refresh list to remove it
		mock.SvcControl_model = FakeObject(GetLocalRec: [])
		Assert(mock.refreshRequired([svc_type: '-', svc_date: Date()]) is: false)
		Assert(mock.refreshRequired([svc_type: '+', svc_date: Date()]) is: false)
		Assert(mock.refreshRequired([svc_type: ' ', svc_date: Date()]) is: false)
		Assert(mock.refreshRequired([svc_type: '-', svc_date: '']) is: false)
		Assert(mock.refreshRequired([svc_type: '+', svc_date: '']) is: false)
		Assert(mock.refreshRequired([svc_type: ' ', svc_date: '']))
		}

	Test_refreshIfCurrentChanged()
		{
		mock = Mock(SvcControl)
		mock.When.refreshIfCurrentChanged().CallThrough()
		mock.When.reselectCurrent([anyArgs:]).CallThrough()
		mock.When.refreshCurrent([anyArgs:]).CallThrough()
		mock.When.checkTreeChanged().Return(false)
		mock.When.On_Refresh().Do({ })
		mock.When.List_Selection([anyArgs:]).Do({ })
		mock.SvcControl_model = FakeObject(GetLocalRec: false)

		mock.refreshIfCurrentChanged()
		mock.Verify.Never().List_Selection([anyArgs:])
		mock.SvcControl_curSelection = #(0)

		_date = Timestamp()
		mock.SvcControl_curSource = class
			{
			Name: 'localList'
			GetRow(@unused)
				{
				return [
					svc_lib: 'Test_lib',
					svc_date: _date,
					svc_type: ' ',
					svc_name: 'test'
					]
				}
			}
		// testing record missing
		mock.refreshIfCurrentChanged()
		mock.Verify.Never().List_Selection([anyArgs:])

		mock.SvcControl_model = FakeObject(GetLocalRec: [lib_modified: _date])
		mock.refreshIfCurrentChanged()
		mock.Verify.Never().List_Selection([anyArgs:])

		mock.SvcControl_model = FakeObject(GetLocalRec: [lib_modified: Timestamp()])
		mock.refreshIfCurrentChanged()
		mock.Verify.List_Selection([anyArgs:])

		// testing new record change
		mock.SvcControl_model = model = Mock()
		model.When.GetLocalRec('Test_lib', 'test', deleted: false).
			Return([lib_modified: Timestamp()])
		mock.SvcControl_local_list = FakeObject(
			GetRow: [
				svc_lib: 'Test_lib',
				svc_date: Timestamp(),
				svc_name: 'test',
				svc_type: '+'
				]
			)
		mock.refreshIfCurrentChanged()
		mock.Verify.Times(2).List_Selection([anyArgs:])
		}

	Test_reselectCurrent()
		{
		mock = Mock(SvcControl)
		mock.When.reselectCurrent([anyArgs:]).CallThrough()

		mockList = Mock()
		mockList.When.SetSelection([anyArgs:]).Do({ })
		mockList.When.Get([anyArgs:]).Return([
				[svc_name: 'Test_Rec1', svc_lib: 'testlib'],
				[svc_name: 'Test_Rec2', svc_lib: 'testlib'],
				[svc_name: 'Test_Rec3', svc_lib: 'testlib'],
				[svc_name: 'Test_Rec1', svc_lib: 'xtestlib'],
			])
		mock.When.FindControl([anyArgs:]).Return(false, mockList)

		mock.reselectCurrent('list_name', [])
		mockList.Verify.Never().Get()

		mock.reselectCurrent('list_name', [])
		mockList.Verify.Get()
		mockList.Verify.Never().SetSelection([anyArgs:])

		mock.reselectCurrent('list_name', [svc_name: 'not_in_list', svc_lib: 'testlib'])
		mockList.Verify.Times(2).Get()
		mockList.Verify.Never().SetSelection([anyARgs:])

		mock.reselectCurrent('list_name', [svc_name: 'Test_Rec1', svc_lib: 'testlib'])
		mockList.Verify.Times(3).Get()
		mockList.Verify.SetSelection(0)

		mock.reselectCurrent('list_name', [svc_name: 'Test_Rec2', svc_lib: 'testlib'])
		mockList.Verify.Times(4).Get()
		mockList.Verify.SetSelection(1)

		mock.reselectCurrent('list_name', [svc_name: 'Test_Rec3', svc_lib: 'testlib'])
		mockList.Verify.Times(5).Get()
		mockList.Verify.SetSelection(2)

		mock.reselectCurrent('list_name', [svc_name: 'Test_Rec1', svc_lib: 'xtestlib'])
		mockList.Verify.Times(6).Get()
		mockList.Verify.SetSelection(3)
		}

	Test_refreshCurrent()
		{
		mock = Mock(SvcControl)
		mock.When.refreshCurrent([anyArgs:]).CallThrough()
		mock.When.On_Refresh().Do({ })
		mock.When.List_Selection([anyArgs:]).Do({ })
		mock.SvcControl_model = FakeObject(GetLocalRec: false)

		mock.refreshCurrent([svc_type: '-'], #(0), Mock())
		mock.Verify.On_Refresh()
		mock.Verify.Never().List_Selection([anyArgs:])

		mock.refreshCurrent([svc_type: ' '], #(0), Mock())
		mock.Verify.Times(2).On_Refresh()
		mock.Verify.Never().List_Selection([anyArgs:])

		mock.refreshCurrent([svc_type: '+'], #(0), Mock())
		mock.Verify.Times(3).On_Refresh()
		mock.Verify.Never().List_Selection([anyArgs:])

		date = Timestamp()
		mock.SvcControl_model = FakeObject(GetLocalRec: [lib_modified: date])
		mock.refreshCurrent([svc_type: ' ', svc_date: date], #(0), Mock())
		mock.Verify.Times(3).On_Refresh()
		mock.Verify.Never().List_Selection([anyArgs:])

		mock.refreshCurrent([svc_type: ' ', svc_date: Timestamp()], #(0), Mock())
		mock.Verify.Times(3).On_Refresh()
		mock.Verify.List_Selection([anyArgs:])
		}

	Test_clearStatus()
		{
		_resultOb = Object(
			errMap:
			Object(svc_all_changes_test1lib_Rec1: 8355839,
			svc_all_changes_test1lib_Rec2: 65280,
			svc_all_changes_test2lib_Rec3: 11196671,
			svc_all_changes_test3lib_Rec4: 8355839,
			svc_all_changes_test4lib_Rec5: 65280,
			svc_all_changes_test4lib_Rec6: 65280,
			svc_all_changes_test2lib_Rec3_Test: 11196671),
			msgMap:
			Object(test1lib_Rec1: "Record has syntax error(s)",
			test2lib_Rec3_Test: "test2lib:Rec3_Test rating is: 4,
				maintain/exceed: 5\n\r\n",
				test3lib_Rec4: "Record has syntax error(s)",
			test2lib_Rec3: "test2lib:Rec3 rating is: 4,
				maintain/exceed: 5\n\r\n"),
			svc_all_changes:
				Object(qualityCheck:
				Object(Object(lib: "test2lib", results: "test2lib:Rec3 rating is: 4,
					maintain/exceed: 5\n", name: "Rec3"),
				Object(lib: "test2lib", results: "test2lib:Rec3_Test rating is: 4,
				maintain/exceed: 5\n", name: "Rec3_Test"))))

		cl = SvcControl
			{
			SvcControl_getWarnings()
				{
				return _resultOb
				}

			SvcControl_getCurTable()
				{
				return "svc_all_changes"
				}

			SvcControl_getSentChanges()
				{
				return Object(Object(lib: "test2lib", name: "Rec3", type: ' '),
					#(lib: "test2lib", name: "Rec3_Test", type: ' '))
				}
			}

		fn = cl.SvcControl_clearStatus

		fn()

		Assert(_resultOb.errMap is: #(svc_all_changes_test1lib_Rec1: 8355839,
			svc_all_changes_test1lib_Rec2: 65280,
			svc_all_changes_test3lib_Rec4: 8355839,
			svc_all_changes_test4lib_Rec5: 65280,
			svc_all_changes_test4lib_Rec6: 65280))

		Assert(_resultOb.msgMap is: #(test1lib_Rec1: "Record has syntax error(s)",
			test3lib_Rec4: "Record has syntax error(s)"))

		Assert(_resultOb.svc_all_changes is: #(qualityCheck: #()))
		}

	Teardown()
		{
		Suneido.Delete(#LibraryTablesOverride)
		Suneido.Delete(#BookTablesOverride)
		super.Teardown()
		}
	}
