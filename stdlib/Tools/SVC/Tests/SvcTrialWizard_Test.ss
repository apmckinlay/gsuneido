// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	cl: SvcTrialWizard
		{
		SvcTrialWizard_trialTags: (test1: 'test1 desc', test2: 'test2 desc')
		}

	Test_main()
		{
		copyAndRestore = .cl.SvcTrialWizard_copyAndRestore
		renameAndDelete = .cl.SvcTrialWizard_renameAndDelete
		lib = .MakeLibrary()
		svcTable = .SvcTable(lib)

		committed = #20251001.1355
		.MakeLibraryRecord(
			[group: 0, parent: 0, num: 1, name: #Folder1],
			[parent: 1, num: 2, lib_modified: Date(), lib_committed: committed,
				text: 'Rec1__modified', lib_before_text: 'Rec1', name: #Rec1__webgui],
			[parent: 1, num: 3, lib_modified: '', lib_committed: committed,
				text: 'Rec2__modified', lib_before_text: '', name: #Rec2__webgui_test2],
			[parent: 1, num: 4, lib_modified: '', lib_committed: committed,
				text: 'Rec2', lib_before_text: '', name: #Rec2__webgui],
			table: lib
			)

		/* Tests on Rec1__webgui */
		// start test1 trail,
		// should restore Rec1__webgui
		// should create Rec1__webgui_test1
		.run(copyAndRestore, 'Rec1__webgui', svcTable, 'test1')

		Assert(svcTable.Get('Rec1__webgui') hasSubset: [parent: 1, group: -1,
			text: 'Rec1', lib_before_text: '',
			lib_committed: committed, lib_modified: ''])
		rec = svcTable.Get('Rec1__webgui_test1')
		Assert(rec hasSubset: [parent: 1, group: -1,
			text: 'Rec1__modified', lib_before_text: '',
			lib_committed: ''])

		// change to test2 trail
		// should permanently delete Rec1__webgui_test1 because it is new
		// should create Rec1__webgui_test2
		.run(renameAndDelete, 'Rec1__webgui_test1', svcTable, 'test2')

		Assert(svcTable.Get('Rec1__webgui_test1') is: false)
		Assert(svcTable.Get('Rec1__webgui_test1', deleted:) is: false)
		rec = svcTable.Get('Rec1__webgui_test2')
		Assert(rec hasSubset: [parent: 1, group: -1,
			text: 'Rec1__modified', lib_before_text: '',
			lib_committed: ''])

		// end test2 trail
		// should permanently delete Rec1__webgui_test2 because it is new
		// should apply changes to Rec1__webgui
		.run(renameAndDelete, 'Rec1__webgui_test2', svcTable, '')
		Assert(svcTable.Get('Rec1__webgui_test2') is: false)
		Assert(svcTable.Get('Rec1__webgui_test2', deleted:) is: false)
		rec = svcTable.Get('Rec1__webgui')
		Assert(rec hasSubset: [parent: 1, group: -1,
			text: 'Rec1__modified', lib_before_text: 'Rec1',
			lib_committed: committed])

		/* Tests on Rec2__webgui */
		// change to test1 trail
		// should mark Rec2__webgui_test2 as deleted
		// should create Rec2__webgui_test1
		.run(renameAndDelete, 'Rec2__webgui_test2', svcTable, 'test1')
		Assert(svcTable.Get('Rec2__webgui_test2') is: false)
		Assert(svcTable.Get('Rec2__webgui_test2', deleted:)
			hasSubset: [parent: 1, group: -2, lib_committed: committed,
				text: 'Rec2__modified', lib_before_text: 'Rec2__modified'])
		Assert(svcTable.Get('Rec2__webgui_test1')
			hasSubset: [parent: 1, group: -1, lib_committed: '',
				text: 'Rec2__modified'])

		// end test1 trail
		// should permanently delete Rec2__webgui_test1 because it is new
		// should apply changes to Rec2__webgui
		.run(renameAndDelete, 'Rec2__webgui_test1', svcTable, '')
		Assert(svcTable.Get('Rec2__webgui_test1') is: false)
		Assert(svcTable.Get('Rec2__webgui_test1', deleted:) is: false)
		Assert(svcTable.Get('Rec2__webgui') hasSubset: [parent: 1, group: -1,
			text: 'Rec2__modified', lib_before_text: 'Rec2',
			lib_committed: committed])

		// start test2 trail
		// should restore Rec2__webgui
		// should restore Rec1__webgui_test2
		.run(copyAndRestore, 'Rec2__webgui', svcTable, 'test2')
		Assert(svcTable.Get('Rec2__webgui') hasSubset: [parent: 1, group: -1,
			text: 'Rec2', lib_before_text: '',
			lib_committed: committed, lib_modified: ''])
		Assert(svcTable.Get('Rec2__webgui_test2') hasSubset: [parent: 1, group: -1,
			text: 'Rec2__modified', lib_before_text: '',
			lib_committed: committed, lib_modified: ''])
		}

	run(fn, name, svcTable, tag)
		{
		Transaction(update:)
			{ |t|
			fn(Object(svc_name: name), svcTable, t, tag)
			}
		}
	}