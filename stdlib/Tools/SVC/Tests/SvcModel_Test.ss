// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_GetLocalRec()
		{
		svcTable = .SvcTable(lib = .MakeLibrary())
		model = .SvcModel(lib)
		svc = model.SvcModel_svc

		addTs = .CommitAdd(svc, svcTable, "One", "one text", "add")
		QueryDo("update " $ lib $ " where name = 'One' set text = 'One modified',
			lib_modified = " $ Display(Timestamp()))
		.assertTextAndDate(model.GetLocalRec(lib, "One"), "One modified", addTs)
		}

	Test_GetMasterRec()
		{
		svcTable = .SvcTable(lib = .MakeLibrary())
		model = .SvcModel(lib)
		svc = model.SvcModel_svc

		.CommitAdd(svc, svcTable, name = .TempName(), "text", "add 1")
		change1Ts = .CommitTextChange(svc, svcTable, name, "change", "change 1")
		change2Ts = .CommitTextChange(svc, svcTable, name, "change again", "change 2")
		delete1Ts = .CommitDelete(svc, svcTable, name, "delete 1")
		add2Ts = .CommitAdd(svc, svcTable, name, "added", "add 2")
		delete2Ts = .CommitDelete(svc, svcTable, name, "delete 2")
		QueryDo("delete " $ lib $" where name = " $ Display(name))

		.assertCommentAndDate(model.GetMasterRec(lib, name, delete?:),
			"delete 2", delete2Ts)
		.assertCommentAndDate(model.GetMasterRec(lib, name, delete1Ts, delete?:),
			"delete 1", delete1Ts)
		.assertTextAndDate(model.GetMasterRec(lib, name, add2Ts),
			"added", add2Ts)
		.assertTextAndDate(model.GetMasterRec(lib, name, change1Ts),
			"change", change1Ts)
		.assertTextAndDate(model.GetMasterRec(lib, name, change2Ts),
			"change again", change2Ts)

		add3Ts = .CommitAdd(svc, svcTable, name, "add again", "add 3")
		.assertTextAndDate(model.GetMasterRec(lib, name, add3Ts), "add again", add3Ts)
		}

	Test_GetPrevMasterRec()
		{
		svcTable = .SvcTable(lib = .MakeLibrary())
		model = .SvcModel(lib)
		svc = model.SvcModel_svc

		add1Ts = .CommitAdd(svc, svcTable, name = .TempName(), "text", "add 1")
		modify1Ts = .CommitTextChange(svc, svcTable, name, "change", "change 1")
		modify2Ts = .CommitTextChange(svc, svcTable, name, "change again", "change 2")
		QueryDo("delete " $ lib $ " where name = " $ Display(name))
		svcTable.SetMaxCommitted(Date.Begin(), force:)

		model = .SvcModel(lib)
		Assert(model.GetPrevMasterRec(lib, name, add1Ts) is: false)
		.assertTextAndDate(
			model.GetPrevMasterRec(lib, name, modify2Ts), "change", modify1Ts)
		.assertTextAndDate(
			model.GetPrevMasterRec(lib, name, modify1Ts), "text", add1Ts)
		}

	Test_GetMasterFromLocal()
		{
		svcTable = .SvcTable(lib = .MakeLibrary())
		model = .SvcModel(lib)
		svc = model.SvcModel_svc
		.CommitAdd(svc, svcTable, name = .TempName(), "add", "add")
		modifyTs = .CommitTextChange(svc, svcTable, name, "changed", "change")
		QueryDo("update " $ lib $ " where name = " $ Display(name) $
			" set text = 'One modified', lib_modified = " $ Display(Timestamp()))

		localChanges = svc.Local_changes(lib)
		Assert(localChanges.FindOne({ it.name is name }).Project(#(type, name))
			is: [type: ' ', :name])
		.assertTextAndDate(model.GetMasterFromLocal(lib, name), "changed", modifyTs)
		}

	Test_GetMergedRec()
		{
		svcTable = .SvcTable(lib = .MakeLibrary())
		model = .SvcModel(lib)
		svc = model.SvcModel_svc
		addTs = .CommitAdd(svc, svcTable, name = .TempName(), "similar", "add")
		.CommitTextChange(svc, svcTable, name, "similar\nmaster", "modify")
		QueryDo("update " $ lib $ " where name = " $ Display(name) $
			" set text = 'similar\nlocal', lib_modified = " $ Display(Timestamp()) $
			", lib_committed = " $
			Display(addTs))
		Assert(model.GetMergedRec(lib, name, model.GetLocalRec(lib, name),
			model.GetMasterRec(lib, name)).merged.Lines()
			is: #("similar", "master", "local"))
		}

	assertTextAndDate(rec, text, lib_committed)
		{
		Assert(rec.lib_committed is: lib_committed)
		Assert(rec.text is: text)
		}

	assertCommentAndDate(rec, comment, lib_committed)
		{
		Assert(rec.comment is: comment)
		Assert(rec.lib_committed is: lib_committed)
		}

	Test_GetPrevDefinition()
		{
		testName = .TempName()
		libs = .makeLibraries(testName)

		model = .SvcModel(lib = .MakeLibrary(), libraries: libs, trialTags: #(trial))

		// If the record is in stdlib, we return prior to the SvcModel.libraries loop
		Assert(model.GetPrevDefinition(testName, 'stdlib') is: false)

		// The first library returned by SvcModel.libraries is essentially stdlib
		// as a result, there can be no previous definitions
		Assert(model.GetPrevDefinition(testName, libs[0]) is: false)
		// All other libraries should successfully return the record from
		// the first library, (our mock stdlib)
		for lib in libs[1 ..]
			{
			rec = model.GetPrevDefinition(testName, lib)
			Assert(rec.table is: libs[0])
			Assert(rec.name is: testName)
			Assert(rec.text is: 'this is record: 0')
			Assert(rec hasMember: 'path')
			}

		// Will not find a previous definition as this record does not exist
		Assert(model.GetPrevDefinition(name = .TempName(), 'stdlib') is: false)
		for lib in libs
			Assert(model.GetPrevDefinition(name, lib) is: false)

		rec = model.GetPrevDefinition(testName $ '__trial', libs[2])
		Assert(rec.table is: libs[2])
		Assert(rec.name is: testName)
		Assert(rec.text is: 'this is record: 2')
			Assert(rec hasMember: 'path')

		rec = model.GetPrevDefinition(testName $ '__trial', libs[4])
		Assert(rec.table is: libs[4])
		Assert(rec.name is: testName)
		Assert(rec.text is: 'this is record: 4')
			Assert(rec hasMember: 'path')
		}

	// Purposely avoiding standard libraries. GetPrevDefinition calls SvcLibrary
	// on each library returned by SvcModel.libraries(). SvcLibrary runs an ensure
	// which will add several columns to the specified table.
	// This will cause SystemChanges.CompareState to fail.
	makeLibraries(testName)
		{
		libs = Object()
		for i in .. 5
			{
			lib = .MakeLibrary([name: testName, text: 'this is record: ' $ i])
			libs.Add(lib)
			Assert(SvcTable.SvcColumns.Intersects?(QueryColumns(lib)) is: false)
			}
		return libs
		}
	}
