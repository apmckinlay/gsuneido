// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	model: false
	Setup()
		{
		.table = .MakeTable('(a,b) key(a)',
			.r1 = [a: 1, b: 'one'],
			.r2 = [a: 2, b: 'two'],
			.r3 = [a: 3, b: 'three'],
			.r4 = [a: 4, b: 'four'])
		.MakeLibraryRecord([name: "Table_" $ .table,
			text: `class { Name: ` $ .table $ ` }`])
		.fktbl = .MakeTable('(c,a) key(c) index (a) in ' $ .table)
		.MakeLibraryRecord([name: "Table_" $ .fktbl,
			text: `class { Name: ` $ .fktbl $ ` }`])
		.TearDownIfTablesNotExist('userselects', 'user_notes')
		}

	Test_main()
		{
		QueryDo('delete ' $ .table $ ' where a >= 5')

		if .model isnt false
			.model.GetCursor().Close()
		.model = AccessModel(Object(.table $ ' sort a'), false)

		.mockAccess = .setupMock(.model, conflictReturns:
			#(false, false, true, false, false, false, true, false, false, true))

		alerts = .doWithAlertRedir()
			{
			.mockAccess.Load_initial_record(#(startLast:))
			Assert(.mockAccess.GetData() is: .r4)

			.mockAccess.Load_initial_record(#())
			Assert(.mockAccess.GetData() is: .r1)

			.check('On_Next', .r2)
			.check('On_Last', .r4)
			.check('On_Prev', .r3)
			.check('On_First', .r1)

			.edit([b: 'ONE'])
			Assert(.mockAccess.GetData() is: [a: 1, b: 'ONE'])

			.check('On_Current_Restore', .r1)

			.check('On_New', [])
			.edit(.r1)

			.mockAccess.On_Current_Restore()
			Assert(.mockAccess.GetData() is: [])

			.edit(.r2)

			// Deleting new rec, goback to r1
			.check('On_Current_Delete', .r1)

			.check('On_First', .r1)

			.check('On_Current_Delete',  .r2) // deleted r1

			.check('On_First', .r2)
			QueryDo('delete ' $ .table $ ' where a = 2') // delete r2
			.check('On_Current_Delete', .r3) // delete a deleted record (just loads next)

			.check('On_Edit', .r3)

			.check('On_Current_Delete', .r4) // delete r3

			.check('On_New', [])
			.edit(.r1)
			Assert(.mockAccess.GetData() is: .r1)
			.check('Save', .r1)
			.check('On_First', .r1)
			.edit(.r1)
			}
		Assert(alerts.Empty?(), msg: 'one')

		QueryDo('update ' $ .table $ ' where a = 1 set b = "ONE"')
		alerts = .doWithAlertRedir()
			{ .check('On_Current_Delete', [a: 1, b: 'ONE']) }
		Assert(alerts[0] hasPrefix: 'Another user has modified this record', msg: 'two')

		alerts = .doWithAlertRedir()
			{
			// fail to delete an updated record
			.mockAccess.Verify.Times(1).AlertError('Current Delete',
				'Another user has modified this record')
			.check('On_Current_Restore', [a: 1, b: 'ONE'])

			r = [a: 1, b: 'uno']
			.edit(r)
			.check('Save', r)

			Assert(Query1(.table $ ' where a = 1') is: r)
			.edit(.r1)
			.check('On_Last', .r4)

			.check('On_New', [])
			.edit(.r3)
			.check('Save', .r3)
			.check('On_New', []) // multiple New
			.edit(.r1)
			}
		Assert(alerts.Empty?(), msg: 'three')

		alerts = .doWithAlertRedir()
			{ .check('On_Last', .r1) }
		Assert(alerts[0] is: 'Duplicate value in field a', msg: 'four')

		alerts = .doWithAlertRedir()
			{
			.edit(.r2)
			.check('Save', .r2)
			.edit(.r1)
			}
		Assert(alerts.Empty?(), msg: 'five')

		alerts = .doWithAlertRedir()
			{ .check('On_Last', .r1) }
		Assert(alerts[0] is: 'Duplicate value in field a', msg: 'six')

		alerts = .doWithAlertRedir()
			{
			.check('On_Current_Restore', .r2)
			.check('On_First', .r1)
			.check('On_Next', .r2)
			.check('On_Next', .r3)
			.check('On_Next', .r4)
			.check('On_Next', .r4)
			}
		Assert(alerts.Empty?(), msg: 'seven')

		// alert message use to come directly from RecordConflict
		alerts = .doWithAlertRedir()
			{
			.edit(r = [a: 4, b: 'FOUR'])
			QueryDo('update ' $ .table $ ' where a = 4 set b = 4444')
			.check('On_First', r)
			}
		Assert(alerts isSize: 0, msg: 'eight')

		alerts = .doWithAlertRedir()
			{
			.check('On_Current_Restore', [a: 4, b: 4444])
			.edit(.r4)
			}
		Assert(alerts.Empty?(), msg: 'nine')

		QueryDo('delete ' $ .table $ ' where a = 4') // delete r4
		alerts = .doWithAlertRedir()
			{ .check('On_First', .r4) }
		Assert(alerts[0] is: "Access: can't get record to update", msg: 'ten')

		alerts = .doWithAlertRedir()
			{ .check('On_Current_Restore', .r4) }
		Assert(alerts[0] is: "The current record has been deleted.", msg: 'eleven')

		alerts = .doWithAlertRedir()
			{
			.check('On_Current_Delete')
			.check('On_First', .r1)
			}
		Assert(alerts.Empty?(), msg: 'twelve')

		QueryOutput(.fktbl, [c: 1, a: 1])
		alerts = .doWithAlertRedir()
			{
			.check('On_Current_Delete', .r1) // delete blocked by foreign key
			}
		Assert(alerts[0] is: 'Record cannot be updated or deleted because it is used',
			msg: 'thirteen')

		alerts = .doWithAlertRedir()
			{
			.edit(r = [a: 11, b: 'one'])
			.check('On_First', r)
			}
		Assert(alerts[0] is: 'Record cannot be updated or deleted because it is used',
			msg: 'fourteen')

		alerts = .doWithAlertRedir()
			{
			.edit(r = [a: 1, b: 'One'])
			.check('On_First', r)
			}
		// Since this came from a save no alert msg
		// use to come directly from RecordConflict
		Assert(alerts isSize: 0, msg: 'fifteen')
		}

	doWithAlertRedir(block)
		{
		alerts = Object()
		alertFn = { |hwnd /*unused*/, message, title /*unused*/, flags /*unused*/|
			alerts.Add(message)
			}
		DoWithAlertRedirected(alertFn)
			{
			(block)()
			}
		return alerts
		}

	setupMock(model, args = #(), conflictReturns = #(false))
		{
		.types = AccessTypes(args)
		mockAccess = Mock(AccessControl)
		nextNum = AccessNextNum(mockAccess, false)
		mockAccess.AccessControl_model = model
		mockAccess.AccessControl_c = model.GetCursor()
		mockAccess.AccessControl_keyquery = model.GetKeyQuery()
		mockAccess.AccessControl_types = .types
		mockAccess.AccessControl_nextNumber = nextNum
		mockAccess.AccessControl_protectField = false
		mockAccess.AccessControl_protect = false
		mockAccess.AccessControl_validField = false
		mockAccess.AccessControl_warningField = false
		mockAccess.AccessControl_customFields = #()
		mockAccess.AccessControl_fields = #(a,b)
		mockAccess.AccessControl_query = .table $ ' sort a'
		mockAccess.AccessControl_saveOnlyLinked = false
		mockAccess.AccessControl_select = false
		mockAccess.AccessControl_base_query = .table
		mockAccess.Vert = Object(Scroll: Object(Hwnd: 0))
		mockAccess.Window = Object(Hwnd: 0)

		r = Mock(RecordControl)
		mockAccess.AccessControl_data = r
		r.When.setup([anyArgs:]).CallThrough()
		r.When.Dirty?([anyArgs:]).CallThrough()
		r.When.NewValue([anyArgs:]).CallThrough()
		r.When.Change([anyArgs:]).CallThrough()
		r.When.Valid([anyArgs:]).CallThrough()
		r.Window = Object(Hwnd: 0)
		r.Change = { |member| r.Change(member) }
		r.setup(false, false, false)
		r.When.HandleFocus().Do({ })
		r.When.Get().CallThrough()
		r.When.Set([anyArgs:]).CallThrough()
		r.When.GetReadOnly([anyArgs:]).CallThrough()

		mockAccess.When.GetData().CallThrough()
		mockAccess.AccessControl_status = FakeObject(Set:, SetDefaultMsg:)
		mockAccess.AccessControl_lock = FakeObject(Trylock:, Unlock: )
		mockAccess.AccessControl_edit_button = FakeObject(Pushed?:, SetEnabled:)
		mockAccess.AccessControl_select_button = FakeObject()
		mockAccess.AccessControl_first_button = FakeObject(Grayed:)
		mockAccess.AccessControl_last_button = FakeObject(Grayed:)
		mockAccess.AccessControl_prev_button = FakeObject(Grayed:)
		mockAccess.AccessControl_next_button = FakeObject(Grayed:)
		mockAccess.AccessControl_loopedAddons = FakeObject(Start:, Stop:)
		mockAccess.AccessControl_menus = FakeObject(UpdateHistory:)
		mockAccess.AccessControl_selectMgr = FakeObject(UsingDefaultFilter?: false)
		mockAccess.Addons = FakeObject(Send:, Collect: Object(''))

		mockAccess.When.Load_initial_record([anyArgs:]).CallThrough()
		mockAccess.When.GetData([anyArgs:]).CallThrough()
		mockAccess.When.On_First([anyArgs:]).CallThrough()
		mockAccess.When.On_Next([anyArgs:]).CallThrough()
		mockAccess.When.On_Last([anyArgs:]).CallThrough()
		mockAccess.When.On_Prev([anyArgs:]).CallThrough()
		mockAccess.When.On_New([anyArgs:]).CallThrough()
		mockAccess.When.On_Current_Restore([anyArgs:]).CallThrough()
		mockAccess.When.On_Current_Delete([anyArgs:]).CallThrough()
		mockAccess.When.Save([anyArgs:]).CallThrough()
		mockAccess.When.FocusFirst([anyArgs:]).Do({ })
		mockAccess.When.SetWhere([anyArgs:]).CallThrough()

		mockAccess.When.RecordConflict?([anyArgs:]).Return(@conflictReturns)
		mockAccess.When.AccessControl_getFocus().Return(0)
		mockAccess.When.AccessControl_setFocus([anyArgs:]).Return(0)
		mockAccess.When.AccessControl_beep([anyArgs:]).Return(0)
		mockAccess.When.AccessControl_deleteOldAttachments().Return(0)
		mockAccess.When.ClearFocus([anyArgs:]).Return(0)
		mockAccess.When.Defer([anyArgs:]).Return(0)
		mockAccess.When.RestoreAttachmentFiles([anyArgs:]).Return(0)
		mockAccess.When.deleteRecordAttachments([anyArgs:]).Return(0)

		return mockAccess
		}

	Test_SetWhere_with_sticky_fields()
		{
		if .model isnt false
			.model.GetCursor().Close()

		.model = AccessModel(Object(.table $ ' sort a'), false)

		.mockAccess = .setupMock(.model, #(stickyFields: #(b)))
		.mockAccess.Load_initial_record(Object())

		.check('On_New', [])
		.edit(r = [a: 5, b: 'sticky'])
		.check('Save', r) // output r5

		Assert(Query1(.table, a: 5) is: [a: 5, b: 'sticky'])

		.check('On_New', [b: 'sticky'])
		.edit([a: 6])
		.check('Save', [a: 6, b: 'sticky']) // output r6
		Assert(Query1(.table, a: 6) is: [a: 6, b: 'sticky'])

		.mockAccess.SetWhere('', quiet:)
		.check('On_New', [b: 'sticky'])
		.edit([a: 7])
		.check('Save', [a: 7, b: 'sticky'])
		Assert(Query1(.table, a: 7) is: [a: 7, b: 'sticky'])

		// Cannot adjust the Query within AccessControl while executing SetWhere
		// Which would change the previous query and current query so it would wipe
		// out the sticky fields, so we are faking it here
		.mockAccess.SetWhere(' where a >= 5 ', quiet:)
		.types.ClearStickyFieldValues()
		.mockAccess.AccessControl_types = .types
		.check('On_New', [])
		.edit([a: 8, b: 'new sticky'])
		.check('Save', [a: 8, b: 'new sticky']) // output r8
		Assert(Query1(.table, a: 8) is: [a: 8, b: 'new sticky'])

		.check('On_New', [b: 'new sticky'])
		.edit([a: 9])
		.check('Save', [a: 9, b: 'new sticky']) // output r9
		Assert(Query1(.table, a: 9) is: [a: 9, b: 'new sticky'])

		QueryDo('delete ' $ .table $ ' where a >= 5')
		}

	check(cmd, data = false)
		{
		if cmd is 'On_New'
			.mockAccess.AccessControl_c = .model.GetCursor()
		else if cmd is 'Save'
			.model.SetKeyQuery(data)

		.mockAccess.AccessControl_keyquery = .model.GetKeyQuery()

		.mockAccess[cmd]()

		if data isnt false
			{
			Assert(.mockAccess.GetData() is: data)
			if cmd in ('On_Next', 'On_New', 'On_Prev', 'On_First', 'On_Last', 'Save')
				.mockAccess.AccessControl_original_record = data
			}
		}

	edit(r)
		{
		.mockAccess.SetEditMode()
		.mockAccess.GetData().Merge(r)
		}

	Test_notAbleToDelete()
		{
		m = AccessControl.AccessControl_notAbleToDelete

		mock = Mock(AccessControl)
		mock.When.AccessControl_readOnlyAccess().Return(true, false, false, false)
		Assert(mock.Eval(m))

		mock.When.Send('Access_AllowDelete').Return(false)
		Assert(mock.Eval(m))

		mock.When.Send('Access_AllowDelete').Return(true)
		mock.When.RecordSet?().CallThrough()
		Assert(mock.Eval(m))

		rule = .TempTableName()
		.MakeLibraryRecord([name: "Rule_" $ rule,
			text: `function () { return "" }`])
		mock.AccessControl_record = [num: 'hello']
		mock.AccessControl_protect = false
		mock.AccessControl_protectField = rule
		mock.AccessControl_newrecord? = true
		mock.When.EditMode?().Return(true)
		mock.When.GetData().Return([num: 'hello'])
		mock.When.GetQuery().Return("")
		mock.AccessControl_model = FakeObject(NotifyObservers: true)
		Assert(mock.Eval(m) is: false)
		}

	Teardown()
		{
		if .model isnt false
			.model.GetCursor().Close()
		super.Teardown()
		}
	}
