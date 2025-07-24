// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_handleDeletedItems_conflicts()
		{
		dbRecs 	= .createTable(4)
		assert_recConflict = function(mock, modRec, dbRec, quiet? = false, never? = false)
			{
			if never?
				mock.BrowseControl_addon.Verify.Never().
					RecordConflict?(modRec, dbRec,
						mock.BrowseControl_query_columns,
						mock.BrowseControl_linkField, :quiet?)
			else
				mock.BrowseControl_addon.Verify.RecordConflict?(modRec, dbRec,
					mock.BrowseControl_query_columns,
					mock.BrowseControl_linkField, :quiet?)
			}

		mock 	= .mockOb(recordChanged?:, recordConflict?:)
		modRecs = Object(.modRec("test0", dirty?:))

		res = .runTran(mock, modRecs)

		assert_recConflict(mock, modRecs[0], dbRecs[0])

		mock.Verify.BrowseControl_setBrowseStatusBar([anyArgs:])
		Assert(res is: false)

		//first rec has no conflicts, so it runs without changing valid/quiet?
		mock 	= .mockOb(recordChanged?:)
		mock.BrowseControl_addon.When.RecordConflict?([anyArgs:]).Return(false, true)
		modRecs = Object(
			.modRec("test0", dirty?:),
			.modRec("test1", dirty?:),
			.modRec("test2", dirty?:),
			.modRec("test3", dirty?:))

		res = .runTran(mock, modRecs, oldrecs: Object())

		assert_recConflict(mock, modRecs[0], dbRecs[0])
		assert_recConflict(mock, modRecs[1], dbRecs[1])
		assert_recConflict(mock, modRecs[2], dbRecs[2], quiet?:)
		assert_recConflict(mock, modRecs[3], dbRecs[3], quiet?:)
		mock.Verify.BrowseControl_setBrowseStatusBar([anyArgs:])
		Assert(res is: false)

		//first rec passes, second one doesn't, seeing how linkfield isnt false
		//it should kick out
		mock = .mockOb(recordChanged?:, linkfield:)
		mock.BrowseControl_addon.
			When.RecordConflict?([anyArgs:]).Return(false, true)

		res = .runTran(mock, modRecs, oldrecs: Object())

		assert_recConflict(mock, modRecs[0], dbRecs[0],)
		assert_recConflict(mock, modRecs[1], dbRecs[1],)
		assert_recConflict(mock, modRecs[2], dbRecs[2], quiet?:, never?:)
		assert_recConflict(mock, modRecs[3], dbRecs[3], quiet?:, never?:)
		mock.Verify.Never().BrowseControl_setBrowseStatusBar([anyArgs:])
		Assert(res is: false)
		}

	Test_handleDeletedItems()
		{
		dbRecs = .createTable(3)

		.runhandleDeletedItems(Object(), deletesSize: 0, oldrecsSize: 0)

		modRecs = Object(
			.modRec("test0"),
			.modRec("test1"),
			.modRec("test2"))
		.runhandleDeletedItems(modRecs, deletesSize: 0, oldrecsSize: 0,
			recordChanged?: false)

		modRecs = Object(
			.modRec("test0", dirty?:),
			.modRec("test1", dirty?:),
			.modRec("test2", dirty?:))
		.runhandleDeletedItems(modRecs, deletesSize: 0, oldrecsSize: 3)

		modRecs = Object(
			.modRec("test0", deleted?:),
			.modRec("test1", deleted?:),
			.modRec("test2", dirty?:))
		.runhandleDeletedItems(modRecs, deletesSize: 2, oldrecsSize: 1)

		// dbRecs[0] and dbRecs[1] should be deleted
		Assert(Query1(.table, keyVal: dbRecs[0].keyVal) is: false)
		Assert(Query1(.table, keyVal: dbRecs[1].keyVal) is: false)
		Assert(Query1(.table, keyVal: dbRecs[2].keyVal) isnt: false)
		}
	runhandleDeletedItems(modRecs, deletesSize, oldrecsSize, recordChanged? = true)
		{
		mock = .mockOb(:recordChanged?)
		oldrecs = Object()
		deletes = Object()

		res = .runTran(mock, modRecs, deletes, oldrecs)

		mock.Verify.Never().BrowseControl_setBrowseStatusBar([anyArgs:])
		Assert(res)
		Assert(deletes isSize: deletesSize)
		Assert(oldrecs isSize: oldrecsSize)
		}

	Test_handleUpdatedItems()
		{
		mock = Mock()
		mock.Addons = Mock()
		fn = BrowseControl.BrowseControl_handleUpdatedItems
		data = Object(
			[type: 'oldrec'],
			[type: 'newRec', Browse_NewRecord:],
			[type: 'dirtyNewRec', Browse_NewRecord:, Browse_RecordDirty:],
			[type: 'deletedRec', listrow_deleted:])
		outputs = Object()
		deletes = Object()
		oldrecs = Object(MockObject(
			Object(
				#(#Update, [type: 'oldrec']),
				#(#Final))
			))
		mock.Eval(fn, data, outputs, deletes, oldrecs)

		Assert(outputs
			is: Object([type: 'dirtyNewRec', Browse_NewRecord:, Browse_RecordDirty:]))
		Assert(deletes
			is: Object(1))
		oldrecs[0].Final() // ensure MockObject was called

		Assert(mock.BrowseControl_nextSaveRec is: [type: 'deletedRec', listrow_deleted:])
		}

	//-------------- Helper Funcs --------------
	createTable(count)
		{
		.table = .MakeTable("(keyVal, test) key(keyVal)")
		dbRecs = Object()
		for(i = 0; i <= count; i++)
			{
			dbRecs[i] = Object(keyVal: Timestamp(), test: "test" $ i)
			QueryOutput(.table, dbRecs[i])
			}
		return dbRecs
		}

	mockOb(recordConflict? = false, linkfield = false)
		{
		mock = Mock(BrowseControl)
		mock.BrowseControl_addon = Mock()
		mock.Addons = Mock()
		mock.BrowseControl_key = #(test)
		mock.BrowseControl_linkField = linkfield
		mock.BrowseControl_list = FakeObject(
			SetSelection: function(unused) {}
			AddHighlight: function(@unused) {}
			Send: function(@args) { return args[0] is 'Browse_AllowDelete'} )
		mock.When.BrowseControl_handleConflictItem([anyArgs:]).CallThrough()
		mock.When.BrowseControl_existingRecordChanged?([anyArgs:]).CallThrough()
		mock.When.BrowseControl_keyQuery([anyArgs:]).CallThrough()
		mock.When.BrowseControl_getBaseQuery([anyArgs:]).CallThrough()
		mock.BrowseControl_addon.When.RecordConflict?([anyArgs:]).Return(recordConflict?)
		mock.When.deleteRecord([anyArgs:]).CallThrough()
		return mock
		}

	runTran(mock, modRecs, deletes = #(), oldrecs = #())
		{
		fn = BrowseControl.BrowseControl_handleDeletedItems
		return Transaction(update:)
			{ |tran| mock.Eval(fn, modRecs, tran, .table, deletes, oldrecs) }
		}

	modRec(orgRec, deleted? = false, dirty? = false)
		{
		modRec = Object(Browse_OriginalRecord: Object(test: orgRec),
			listrow_deleted: deleted?)
		if dirty?
			modRec.Add(true at: "Browse_RecordDirty")
		return modRec
		}
	}
