// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
SvcTests
	{
	Test_local_changes()
		{
		// local changes
		.SpyOn(SvcCore.SvcCore_svcHooks).Return('')
		alertCalls = .SpyOn(Alert).Return('').CallLogs()
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		records = [
			[parent: 0, name: 'One', text: 'one text'],
			[parent: 0, name: 'Two', text: 'two text'],
			[parent: 0, name: 'Three', text: 'three text'],
			[parent: 0, name: 'Four', text: 'four text'],
			[parent: 0, name: 'Seven', text: 'seven text']]
		records.Each(svcTable.Output)
		result = svc.Local_changes(lib)
		Assert(result.Map({ it.type $ ' ' $ it.name })
			is: #('+ Four', '+ One', '+ Seven', '+ Three', '+ Two'))

		for x in records
			Assert(svc.Put(svcTable, x.name, 'default', 'new') isDate: true)
		.AssertSvcEmpty(svc, lib)

		// record modified
		QueryDo("update " $ lib $ " where name = 'One' set text = 'One modified',
			lib_modified = " $ Display(Timestamp()))
		Assert(svc.Local_changes(lib)[0].Project(#(type, name))
			is: #(type: ' ', name: 'One'))
		Assert(svc.Put(svcTable, 'One', 'default', 'updated') isDate: true)
		.AssertSvcEmpty(svc, lib)

		// record deleted
		svcTable.StageDelete('Four')
		Assert(svc.Local_changes(lib)[0].Project(#(type, name))
			is: #(type: '-', name: 'Four'))
		svc.Remove(svcTable, 'Four', 'default', 'removed')
		.AssertSvcEmpty(svc, lib)

		// restored changes
		QueryDo('update ' $ lib $ ' where name = #Two
			set lib_before_text = text, text = "two modified",
			lib_before_path = "' $ lib $ '/two_rec_path",
			lib_modified = ' $ Display(Timestamp()))
		svc.Restore(lib, 'Two')
		rec = Query1(lib, name: 'Two')
		Assert(rec.text is: 'two text')
		Assert(rec.lib_before_path is: '')
		// Ensuring folder path is correct, despite not existing prior to restore
		folder = Query1(lib, name: #two_rec_path)
		Assert(rec.parent is: folder.num)
		Assert(folder.parent is: 0)
		// Discrepancy is from us forcing lib_before_path to be something other than
		// it's real path (""). Discrepancies like this will be flagged during the restore
		Assert(alertCalls[0].msg startsWith: 'Unexpected discrepancy detected:\n    path')
		.AssertSvcEmpty(svc, lib)

		// restored changes (text is modified)
		// Ensuring the path / parent are not wrongfully set to 0 (root)
		QueryDo('update ' $ lib $ ' where name = #Two
			set lib_before_text = text, text = "two modified",
			lib_modified = ' $ Display(Timestamp()))
		svc.Restore(lib, 'Two')
		rec = Query1(lib, name: 'Two')
		Assert(rec.text is: 'two text')
		Assert(rec.parent is: folder.num)
		Assert(alertCalls[1].msg startsWith: 'Unexpected discrepancy detected:\n    path')
		.AssertSvcEmpty(svc, lib)

		// restore deletes
		QueryDo('update ' $ lib $ ' where name = #Three
			set	lib_before_path = "' $ lib $ '/extra_path/three_rec_path"')
		svcTable.StageDelete('Three')
		svc.Restore(lib, 'Three')
		rec = Query1(lib, name: 'Three')
		Assert(rec.text is: 'three text')
		Assert(rec.lib_before_path is: '')
		// Ensuring folder path is correct, despite not existing prior to restore
		folder = Query1(lib, name: #three_rec_path)
		Assert(rec.parent is: folder.num)
		folderParent = Query1(lib, name: #extra_path)
		Assert(folder.parent is: folderParent.num)
		Assert(folderParent.parent is: 0)
		Assert(alertCalls[2].msg startsWith: 'Unexpected discrepancy detected:\n    path')
		.AssertSvcEmpty(svc, lib)

		// conflicting changes - see also Test_conflicts()
		masterModified = .CommitTextChange(
			svc, svcTable, 'One', 'one modified', 'master modified')
		.OutstandingConflict(svc, svcTable, 'One', 'one modified by local')

		localModified = Query1(lib, name: 'One', group: -1).lib_modified
		Assert(svc.Conflicts(lib) isSize: 1)
		Assert(svc.Conflicts(lib)[0].
			Project(#(name, localModified, masterModified, localType, masterType))
			is: Object(name: 'One', :localModified, :masterModified,
				localType: ' ', masterType: ' '))
		}

	Test_master_changes()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// add + modify + modify + delete + add + delete
		.CommitAdd(svc, svcTable, "One", "one text", "add 1")
		.CommitTextChange(svc, svcTable, "One", "one changed", "change 1")
		.CommitTextChange(svc, svcTable, "One", "one changed again", "change 2")
		.CommitDelete(svc, svcTable, "One", "delete 1")
		.CommitAdd(svc, svcTable, "One", "one added", "add 2")
		.CommitDelete(svc, svcTable, "One", "delete 2")

		// Clear local table to properly simulate getting the master changes
		QueryDo('delete ' $ lib)
		Assert(svc.Master_changes(lib) isSize: 6)
		Assert(svc.Master_changes(lib).Map!({ it.type $ ' ' $ it.name })
			is: #('+ One', '  One', '  One', '- One', '+ One', '- One'))
		}

	Test_conflicts()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())

		add1Ts = .CommitAdd(svc, svcTable, "One", "one", "add")
		modify1Ts = .CommitTextChange(svc, svcTable, "One", "one master 1", "change 1")
		modify2Ts = .CommitTextChange(svc, svcTable, "One", "one master 2", "change 2")
		deleteTs = .CommitDelete(svc, svcTable, "One", "delete")
		masterModified = .CommitAdd(svc, svcTable, "One", "one master 3", "add")
		.OutstandingConflict(svc, svcTable, 'One', 'local', add1Ts)

		localModified = Query1(lib, name: 'One', group: -1).lib_modified
		Assert(svc.Conflicts(lib) isSize: 1)
		Assert(svc.Conflicts(lib)[0].
			Project(#(name, localModified, masterModified, localType, masterType))
			is: Object(name: 'One', :localModified, :masterModified,
				localType: ' ', masterType: '+'))
		Assert(svc.Conflicts(lib)[0].sends isSize: 4)
		Assert(svc.Conflicts(lib)[0].sends.Map({
			it.Project(#(name, masterType, masterModified)) })
			is: Object(
				Object(name: 'One', masterType: ' ', masterModified: modify1Ts),
				Object(name: 'One', masterType: ' ', masterModified: modify2Ts),
				Object(name: 'One', masterType: '-', masterModified: deleteTs),
				Object(name: 'One', masterType: '+', :masterModified)))
		}

	Test_UpdateLibrary_Changes()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		// Record is created locally, addition is committed to svc
		.CommitAdd(svc, svcTable, 'TestChanges', 'initial text', 'add')
		// Record is modified locally, changes are committed to svc
		.CommitTextChange(svc, svcTable, 'TestChanges', 'text changed x 1', 'change 1')
		// Record is modified locally again, changes are committed to svc
		.CommitTextChange(svc, svcTable, 'TestChanges', 'text changed x 2', 'change 2')
		// Library is cleared, required to simulate getting the changes for the first time
		QueryDo('delete ' $ lib)

		// Svc should list all of the changes, (add, modification 1, modification 2)
		Assert(svc.Master_changes(lib) isSize: 3)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 1)
		.AssertSvcEmpty(svc, lib)
		Assert(Query1(lib, name: 'TestChanges').text is: 'text changed x 2')
		}

	Test_UpdateLibrary_Conflict()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		// Record is created locally, addition is committed to svc
		.CommitAdd(svc, svcTable, 'TestConflict', 'inital text', 'add 1')
		// Record is modified locally, changes are committed to svc
		.CommitTextChange(svc, svcTable, 'TestConflict', 'text changed x 1', 'change 1')
		// Record is modified locally again, changes are committed to svc
		.CommitTextChange(svc, svcTable, 'TestConflict', 'text changed x 2', 'change 2')
		// A conflict occurs
		.OutstandingConflict(svc, svcTable, 'TestConflict', 'local text')
		Assert(svc.Conflicts(lib) isSize: 1)

		model = .SvcModel(lib)
		model.SetSettings([svc_server: false, svc_local?: false])
		model.SetTable(lib)
		model.MoveConflict('TestConflict', lib, merge?:)
		Assert(model.UpdateLibrary(lib) is: 1)

		rec = Query1(lib, name: 'TestConflict')
		Assert(rec.text is: 'text changed x 2')
		Assert(rec.lib_invalid_text is: 'text changed x 2\r\nlocal text')
		Assert(rec.lib_current_text is: 'text changed x 2\r\nlocal text')
		Assert(svcTable.GetMaxCommitted() is: rec.lib_committed)
		}

	Test_UpdateLibrary_Conflict_suffix()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// Testing conflict with a new record
		ts = .CommitAdd(svc, svcTable, 'Test0', 'text', 'adding').Minus(minutes: 1)
		.OutstandingConflict(svc, svcTable, 'Test0', 'conflict Test0', '')
		Assert(svc.Conflicts(lib) isSize: 1)

		// Testing conflicts with a record with one modification
		.CommitAdd(svc, svcTable, 'Test1', 'text', 'adding')
		.CommitTextChange(svc, svcTable, 'Test1', 'changed', 'change 1')
		.OutstandingConflict(svc, svcTable, 'Test1', 'conflict Test1', ts)
		Assert(svc.Conflicts(lib) isSize: 2)

		// Testing conflicts with a record multiple modifications
		.CommitAdd(svc, svcTable, 'Test2', 'text', 'adding')
		.CommitTextChange(svc, svcTable, 'Test2', 'changed', 'change 1')
		.CommitTextChange(svc, svcTable, 'Test2', 'changed', 'change 2')
		.OutstandingConflict(svc, svcTable, 'Test2', 'conflict Test2', ts)
		Assert(svc.Conflicts(lib) isSize: 3)
		}

	Test_UpdateLibrary_AddDelete()
		{
		startTime = Timestamp()
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())

		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		// Record is created locally, addition is committed to svc
		.CommitAdd(svc, svcTable, '3', 'initial add', 'add')
		// Record is modified locally, changes are committed to svc
		.CommitTextChange(svc, svcTable, '3', 'changed', 'change 1')
		// Record is deleted, deletion is committed to svc
		.CommitDelete(svc, svcTable, '3', 'deleted')
		// Have to reset the libCommitted record to see these changes
		svcTable.SetMaxCommitted(Date.Begin(), force:)

		// lib ends up being empty after the delete so there is no max lib_committed
		// as a result, all changes committed to the master table is present
		Assert(svc.Master_changes(lib) isSize: 3)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 1)
		Assert(Query1(lib, name: '3') is: false)
		Assert(Query1(svcTable.NameQuery('3', deleted:)) is: false)

		Assert(svcTable.GetMaxCommitted() greaterThan: startTime)
		}

	Test_UpdateLibrary_Delete()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		.OutstandingDelete(svc, svcTable, '4', '4 text 1', 'outstanding change')

		Assert(Query1(lib, name: '4') isnt: false)
		Assert(svc.Master_changes(lib) isSize: 1)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 1)
		Assert(Query1(lib, name: '4') is: false)
		}

	Test_UpdateLibrary_CommitDeleteAdd()
		{
		startTime = Timestamp()
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		.CommitAdd(svc, svcTable, '5', '5 text', 'add')
		Assert(svc.Master_changes(lib) isSize: 0)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 0)

		QueryOutput(lib $ '_master',
			[lib_committed: Timestamp(), comment: "remove", type: "-", name: "5"])
		QueryOutput(lib $ '_master',
			[text: "5 text", lib_committed: Timestamp(), comment: "add", type: "+",
				name: "5"])
		Assert(svc.Master_changes(lib) isSize: 2)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 1)

		Assert(svcTable.GetMaxCommitted() greaterThan: startTime)
		}

	Test_UpdateLibrary_CommitDeleteAddDelete()
		{
		startTime = Timestamp()
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		.CommitAdd(svc, svcTable, '6', '6 text', 'add')
		Assert(svc.Master_changes(lib) isSize: 0)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 0)

		QueryOutput(lib $ '_master',
			[lib_committed: Timestamp(), comment: "remove", type: "-", name: "6"])
		QueryOutput(lib $ '_master',
			[text: "6 text", lib_committed: Timestamp(), comment: "add", type: "+",
				name: "6"])
		QueryOutput(lib $ '_master',
			[lib_committed: Timestamp(), comment: "deleted", type: "-", name: "6"])
		Assert(svc.Master_changes(lib) isSize: 3)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 1)

		Assert(svcTable.GetMaxCommitted() greaterThan: startTime)
		}

	Test_UpdateLibrary_DeleteAddDelete()
		{
		startTime = Timestamp()
		svcTable = .SvcTable(lib = .MakeLibrary())
		SvcCore.EnsureMaster(lib $ '_master', 'lib')
		svc = .Svc()
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		// Local out of sync, never got the original addition for record 7
		QueryOutput(lib $ '_master',
			[lib_committed: Timestamp(), comment: "deleted", type: "-", name: "7"])
		QueryOutput(lib $ '_master',
			[text: "7 text", lib_committed: Timestamp(), comment: "add", type: "+",
				name: "7"])
		QueryOutput(lib $ '_master',
			[lib_committed: Timestamp(), comment: "deleted", type: "-", name: "7"])
		Assert(svc.Master_changes(lib) isSize: 3)
		Assert(svc.UpdateLibrary(svc.Master_changes(lib)) is: 1)

		Assert(svcTable.GetMaxCommitted() greaterThan: startTime)
		}

	Test_UpdateLibrary_Multi()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		// prevent messages from going to the IDE
		.SpyOn(SvcTable.Publish).Return(0)

		// Record is output and modified. We put it into a conflict state later
		.CommitAdd(svc, svcTable, 'Rec2', 'conflict', 'added')
		.CommitTextChange(svc, svcTable, 'Rec2', 'conflict changed', 'change 1')

		// Outstanding Delete (local record exists, delete exists in svc)
		ts = .CommitAdd(svc, svcTable, 'Rec0', 'ToBeDeleted', 'added')
		svc.Remove(svcTable, 'Rec0', 'default', 'record is deleted')

		// Record is output, then deleted
		.CommitAdd(svc, svcTable, 'Rec1', 'Added then deleted', 'added')
		.CommitTextChange(svc, svcTable, 'Rec1', 'text changed', 'change 1')
		.CommitDelete(svc, svcTable, 'Rec1', 'record is deleted')

		// Put record into conflict state
		.CommitTextChange(svc, svcTable, 'Rec2', 'conflict changed again', 'change 2')
		.OutstandingConflict(svc, svcTable, 'Rec2', 'local changes')

		// Record is output, then modified, then deleted locally
		.CommitAdd(svc, svcTable, 'Rec3', 'text', 'add')
		.CommitTextChange(svc, svcTable, 'Rec3', 'text changed', 'change 1')
		.CommitTextChange(svc, svcTable, 'Rec3', 'text changed again', 'change 2')

		// Delete Rec1 so we can properly simulate getting a delete from svc
		QueryDo('delete ' $ lib $ ' where name in ("Rec1")')

		// Setting lib_committed to be earlier value,
		// This will dictate which changes are seen as outstanding
		QueryDo('update ' $ lib $ ' set lib_committed = ' $ Display(ts))

		model = .SvcModel(lib)
		// Should have 1 conflict and 7 changes
		Assert(model.Conflicts isSize: 1)
		Assert(model.MasterChanges isSize: 7)
		// Moving the conflict will produce two new records in master list
		// 1. The original change record (Rec2: ' ')
		// 2. The merged record (Rec2: '#')
		model.MoveConflict('Rec2', lib, merge?:)
		// Conflicts should now be empty, master changes will now have the two new records
		Assert(model.Conflicts isSize: 0)
		Assert(model.MasterChanges isSize: 9)
		Assert(model.UpdateLibrary(lib) is: 4)

		rec = Query1(lib, name: 'Rec3')
		Assert(rec.text is: 'text changed again')
		Assert(rec.lib_invalid_text is: '')
		Assert(rec.lib_current_text is: 'text changed again')

		rec = Query1(lib, name: 'Rec2')
		Assert(rec.text is: 'conflict changed again')
		Assert(rec.lib_invalid_text is: 'conflict changed again\r\nlocal changes')
		Assert(rec.lib_current_text is: 'conflict changed again\r\nlocal changes')
		}

	Test_list_all_changes()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		name = 'TestRecord'
		dir = 'testDir'
		QueryOutput(lib, [num: 8, :name, text: '//This is the record', group: -1])

		svccore = Mock(SvcCore)
		svccore.When.ensureDir([anyArgs:]).Do({ })
		svccore.When.export([anyArgs:]).Do({ })
		svccore.When.addToDeletes([anyArgs:]).Do({ })
		svccore.When.masterTableList([anyArgs:]).Return(Object(lib $ '_master'))
		svccore.When.ListAllChanges([anyArgs:]).CallThrough()

		Assert(svc.Put(svcTable, name, 'default', 'new') isDate: true)
		Assert(svccore.ListAllChanges(date = Date().Plus(hours: 1)) is: #())
		svccore.Verify.Never().ensureDir([anyArgs:])

		Assert(svccore.ListAllChanges(date, :dir) is: #())
		svccore.Verify.ensureDir(dir)

		date = Date().Minus(hours: 1)
		Assert(svccore.ListAllChanges(date, :dir)[lib] is: [name])
		svccore.Verify.export([anyArgs:])
		svccore.Verify.Never().addToDeletes([anyArgs:])
		svccore.Verify.Never().bookInfo([anyArgs:])

		Assert(svc.Remove(svcTable, name, 'default', 'delete') isDate: true)
		Assert(svccore.ListAllChanges(date, :dir)[lib] is: ['-' $ name])
		svccore.Verify.addToDeletes(dir, lib, name)
		svccore.Verify.export([anyArgs:])
		svccore.Verify.Never().bookInfo([anyArgs:])
		}

	Test_MissingTest?()
		{
		lib = .MakeLibrary()
		mock = Mock()
		mock.When.Get(lib, 'Abc_Test').Return([])
		Assert(mock.Eval(Svc.MissingTest?, lib, 'Abc') is: false)

		mock.When.Get(lib, 'Bcd_Test').Return(false)
		mock.When.Get(lib, 'BcdTest').Return([])
		Assert(mock.Eval(Svc.MissingTest?, lib, 'Bcd') is: false)

		mock.When.Get(lib, 'Cde_Test').Return([])
		Assert(mock.Eval(Svc.MissingTest?, lib, 'Cde?') is: false)

		mock.When.Get(lib, 'Efg_Test').Return(false)
		mock.When.Get(lib, 'EfgTest').Return([])
		Assert(mock.Eval(Svc.MissingTest?, lib, 'Efg?') is: false)

		mock.When.Get(lib, 'Fgh_Test').Return(false)
		mock.When.Get(lib, 'FghTest').Return(false)
		Assert(mock.Eval(Svc.MissingTest?, lib, 'Fgh'))

		mock.When.Get(lib, 'Ghi_Test').Return(false)
		mock.When.Get(lib, 'GhiTest').Return(false)
		Assert(mock.Eval(Svc.MissingTest?, lib, 'Ghi?'))
		}

	Test_hashes()
		{
		.SpyOn(SvcCore.SvcCore_svcHooks).Return('')
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		firstHash = Svc.Hash('1')
		secondHash = Svc.Hash('2')
		thirdHash = Svc.Hash('3')
		fourthHash = Svc.Hash('4')

		.CommitAdd(svc, svcTable, name = 'Hash1', '1', 'add')
		svcRec = svc.Get(lib, name)
		Assert(svcRec.lib_before_hash is: '')
		Assert(svcRec.lib_after_hash is: firstHash)
		Assert(Query1(lib, :name).lib_before_hash is: firstHash)

		.CommitTextChange(svc, svcTable, name, '2', 'text change')
		svcRec = svc.Get(lib, name)
		Assert(svcRec.lib_before_hash is: firstHash)
		Assert(svcRec.lib_after_hash is: secondHash)
		Assert(Query1(lib, :name).lib_before_hash is: secondHash)

		// rename
		.CommitDelete(svc, svcTable, name, 'Renamed Hash1 to Hash2')
		.CommitAdd(svc, svcTable, name = 'Hash2', '3', 'Renamed Hash1 to Hash2',
			lib_before_hash: secondHash)
		svcRec = svc.Get(lib, name)
		Assert(svcRec.lib_before_hash is: secondHash)
		Assert(svcRec.lib_after_hash is: thirdHash)
		Assert(Query1(lib, :name).lib_before_hash is: thirdHash)

		// Simulate lib_before_hash not being set during a text change
		// Svc will use: lib_before_text to recalculate lib_before_text
		QueryDo('update ' $ lib $ ' where name is "' $ name $ '"
			set lib_before_hash = "", lib_before_text = text, text = "4"')
		Assert(svc.Put(svcTable, name, 'default', 'fake change') isDate: true)
		svcRec = svc.Get(lib, name)
		Assert(svcRec.lib_before_hash is: thirdHash)
		Assert(svcRec.lib_after_hash is: fourthHash)
		Assert(Query1(lib, :name).lib_before_hash is: fourthHash)

		// Simulate lib_before_hash not being set during a record move
		// Svc will use: text to recalculate lib_before_text
		QueryDo('update ' $ lib $ ' where name is "' $ name $ '"
			set lib_before_hash = "", lib_modified = ' $ Display(Timestamp()))
		Assert(svc.Put(svcTable, name, 'default', 'fake move') isDate: true)
		svcRec = svc.Get(lib, name)
		Assert(svcRec.lib_before_hash is: fourthHash)
		Assert(svcRec.lib_after_hash is: fourthHash)
		Assert(Query1(lib, :name).lib_before_hash is: fourthHash)
		}

	Test_Library?()
		{
		Assert(.Svc().Library?(.MakeLibrary()))
		Assert(.Svc().Library?(.MakeBook()) is: false)
		}

	Test_CheckCommitted()
		{
		svcTable = .SvcTable(.MakeLibrary())
		m = .Svc().CheckCommitted
		alertCalls = .SpyOn(Alert).Return('').CallLogs()

		m(svcTable, 'TestRec', false, false)
		Assert(alertCalls isSize: 0)

		local = [lib_committed: Date(), text: 'Order: 0\r\n\r\nTesting Compare']
		m(svcTable, 'TestRec', local, false)
		Assert(alertCalls isSize: 0)

		master = local.Copy()
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls isSize: 0)

		master.lib_committed = master.lib_committed.Plus(minutes: 1)
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls[0].msg
			startsWith: 'Unexpected discrepancy detected:\n    lib_committed')

		master.text = 'Different'
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls[1].msg
			startsWith: 'Unexpected discrepancy detected:' $
				'\n    lib_committed\n    text')

		// SvcTable is a SvcLibrary, so order is treated as part of text
		master.text = 'Order: 0\n\nTesting Compare'
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls[2].msg
			startsWith: 'Unexpected discrepancy detected:' $
				'\n    lib_committed\n    text')

		master.path = 'root/folder/etc'
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls[3].msg
			startsWith: 'Unexpected discrepancy detected:' $
				'\n    lib_committed\n    path\n    text')

		// SvcTable is now a SvcBook, so Order and text are compared separately
		// Should match despite text separating order and text with \n (vs \r\n)
		m(svcTable = .SvcTable(.MakeBook()), 'TestRec', local, master)
		Assert(alertCalls[4].msg
			startsWith: 'Unexpected discrepancy detected:' $
				'\n    lib_committed\n    path')

		master.text = 'Order: 1\r\n\r\nTesting Compare'
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls[5].msg
			startsWith: 'Unexpected discrepancy detected:' $
				'\n    lib_committed\n    path\n    order')

		master.path = ''
		master.text = 'Order: 1\n\nTesting Compare'
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls[6].msg
			startsWith: 'Unexpected discrepancy detected:' $
				'\n    lib_committed\n    order')

		// Book records should match regardless of the newlines used to append order
		master.text = 'Order: 0\n\nTesting Compare'
		master.lib_committed = local.lib_committed
		Assert(alertCalls isSize: 7)
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls isSize: 7)

		master.text = 'Order: 0\r\n\r\nTesting Compare'
		master.lib_committed = local.lib_committed
		Assert(alertCalls isSize: 7)
		m(svcTable, 'TestRec', local, master)
		Assert(alertCalls isSize: 7)
		}

	Test_SendLocalChanges()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		.SpyOn(SvcCore.SvcCore_svcHooks).Return('')
		.prevMaxCommitted = .prevLibCommitted = false

		Assert(svc.SendLocalChanges([], 'no changes', #default))

		// Put returns false, meaning someone else sent changes,
		// or couldn't find the record
		changes = [
			[:lib, type: '+', name: #TestRecord0]
			]
		Assert(svc.SendLocalChanges(changes, #Test, #default) is: false)

		// Put returns true. One change is sent successfully
		Assert(QueryAll(lib $ ' where group isnt -3') isSize: 0)
		.testSend(changes, svc, svcTable, 'Test Send 1')

		// Put returns true. Multiple changes are sent successfully
		changes = [
			[:lib, type: ' ', name: #TestRecord0],
			[:lib, type: '+', name: #TestRecord1],
			[:lib, type: '+', name: #TestRecord2],
			[:lib, type: '+', name: #TestRecord3],
			[:lib, type: '+', name: #TestRecord4],
			[:lib, type: '+', name: #TestRecord5],
			]
		.testSend(changes, svc, svcTable, 'Test Send 2')

		// Put returns true. Multiple changes are sent successfully
		changes = [
			[:lib, type: ' ', name: #TestRecord0],
			[:lib, type: ' ', name: #TestRecord1],
			[:lib, type: ' ', name: #TestRecord2],
			[:lib, type: '+', name: #TestRecord6],
			[:lib, type: '-', name: #TestRecord3],
			[:lib, type: '-', name: #TestRecord4],
			[:lib, type: '-', name: #TestRecord5],
			]
		.testSend(changes, svc, svcTable, 'Test Send 3')

		// Only deletes are sent, ensure updateMaxLibCommitted updates accordingly
		// (see: assertLibCommittedDates)
		changes = [
			[:lib, type: '-', name: #TestRecord0],
			[:lib, type: '-', name: #TestRecord1],
			[:lib, type: '-', name: #TestRecord2],
			]
		.testSend(changes, svc, svcTable, 'Test Send 4')

		changes = [
			[:lib, type: '-', name: 'TestRecord6'],
			]
		.testSend(changes, svc, svcTable, 'Test Send 5')
		}

	testSend(changes, svc, svcTable, desc, id = 'default')
		{
		send = Timestamp()
		changes.Each()
			{
			if it.type is '+'
				.outputLocalRec(it.name, it.lib)
			else if it.type is '-'
				svcTable.StageDelete(it.name)
			}
		Assert(svc.SendLocalChanges(changes, desc, id))
		masterTable = svcTable.Table() $ '_master'
		changes.Each()
			{
			// One change for this record should exist after the send date
			// this record should match the sent data
			Assert(Query1(masterTable $
				' where lib_committed > ' $ Display(send),
				name: it.name, comment: desc, type: it.type)
				isnt: false)
			Assert(QueryColumns(masterTable) has: 'svc_library')
			if it.type is '-'
				{
				Assert(Query1(svcTable.NameQuery(it.name)) is: false)
				Assert(Query1(svcTable.NameQuery(it.name, deleted:)) is: false)
				}
			}
		.assertLibCommittedDates(svcTable)
		}

	outputLocalRec(name, lib)
		{
		rec = [:name, text: `function () { Print("` $ name $ `") }`]
		OutputLibraryRecord(lib, rec)
		}

	assertLibCommittedDates(svcTable)
		{
		newMaxCommitted = svcTable.GetMaxCommitted()
		Assert(newMaxCommitted isnt: .prevMaxCommitted)
		.prevMaxCommitted = newMaxCommitted

		newLibCommitted = svcTable.GetMaxCommitted()
		Assert(newLibCommitted isnt: .prevLibCommitted)
		.prevLibCommitted = newMaxCommitted
		}

	Test_buildMsg()
		{
		fn = Svc.Svc_buildMsg
		Assert(fn(changeType: '+', lib: 'fakelib', name: 'name', prefix: '<<<')
			is: '<<<+fakelib:name')
		Assert(fn(changeType: '-', lib: 'fakelib', name: 'name', failed?:)
			is: '!!! -fakelib:name FAILED !!!')
		}
	}