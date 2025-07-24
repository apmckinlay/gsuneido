// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	cl: LibTreeModel
			{
			New(.returnLibs)
				{
				super()
				}
			LibTreeModel_initCache()
				{
				mock = Mock()
				mock.When.Reset().Do({ })
				return mock
				}
			Libs()
				{
				return .returnLibs
				}
			}
	treeModel(tables)
		{
		return new .cl(tables)
		}

	Test_renameRecord()
		{
		// Setup
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		.AssertSvcEmpty(svc, lib)

		// Output records and get their nums
		recs = [
			[name: 'One',   text: 'one text'],
			[name: 'Two',   text: 'two text'],
			[name: 'Three', text: 'three text'],
			[name: 'Four',  text: 'four text']].
			Each(svcTable.Output)
		oneNum = recs[0].num
		twoNum = recs[1].num
		threeNum = recs[2].num
		fourNum  = recs[3].num

		// Assert the new records are seen as outstanding. Test sending to SVC
		result = svc.Local_changes(lib)
		Assert(result.Map({ it.type $ ' ' $ it.name })
			is: #('+ Four', '+ One', '+ Three', '+ Two'))
		for x in recs
			Assert(svc.Put(svcTable, x.name, 'default', 'new') isDate:)
		.AssertSvcEmpty(svc, lib)

		// rename record 1
		treeModel = .treeModel([lib])
		Assert(recOne = treeModel.Get(treeModel.MangleNum(lib, oneNum)) isnt: #())
		treeModel.Rename(recOne, 'xOne')
		// make sure there's a delete and an add
		Assert(svc.Local_changes(lib).Map({ it.type $ ' ' $ it.name })
			is: #('+ xOne', '- One'))
		recOneDel = svcTable.Get('One', deleted:)
		Assert(recOneDel.lib_committed isnt: '')
		Assert(recOneDel.lib_before_text is: 'one text')

		// rename record 2 to record 1
		Assert(recTwo = treeModel.Get(treeModel.MangleNum(lib, twoNum)) isnt: #())
		treeModel.Rename(recTwo, 'One')
		Assert(recTwo = svcTable.Get('One') isnt: false)
		Assert(recTwo.lib_before_text is: 'one text')

		// make sure there's a del rec 2, add ren. rec 1, change for rec 1
		Assert(svc.Local_changes(lib).Map({ it.type $ ' ' $ it.name })
			is: #('  One', '+ xOne', '- Two'))
		recOneDel = svcTable.Get('One', deleted:)
		recTwoDel = svcTable.Get('Two', deleted:)
		Assert(recOneDel is: false)
		Assert(recTwoDel.lib_committed isnt: '')
		Assert(recTwoDel.lib_before_text is: 'two text')

		Assert(svc.Put(svcTable, 'xOne', 'default', 'new') isDate:)
		Assert(svc.Put(svcTable, 'One', 'default', 'updated') isDate:)
		svc.Remove(svcTable, 'Two', 'default' 'removed')
		.AssertSvcEmpty(svc, lib)

		// rename Three to xThree, Four to xFour, xThree to Four, xFour to Three
		Assert(recThree = treeModel.Get(treeModel.MangleNum(lib, threeNum)) isnt: #())
		treeModel.Rename(recThree, 'xThree')
		Assert(svcTable.Get('xThree') isnt: false)

		Assert(recFour = treeModel.Get(treeModel.MangleNum(lib, fourNum)) isnt: #())
		treeModel.Rename(recFour, 'xFour')
		Assert(svcTable.Get('xFour') isnt: false)

		recThreeDel = svcTable.Get('Three', deleted:)
		recFourDel = svcTable.Get('Four', deleted:)
		Assert(recThreeDel.lib_committed isnt: '')
		Assert(recThreeDel.lib_before_text is: 'three text')
		Assert(recFourDel.lib_committed isnt: '')
		Assert(recFourDel.lib_before_text is: 'four text')

		Assert(recThree = treeModel.Get(treeModel.MangleNum(lib, threeNum)) isnt: #())
		treeModel.Rename(recThree, 'Four')
		Assert(recThree = svcTable.Get('Four') isnt: false)
		Assert(recThree.lib_before_text is: 'four text')

		Assert(recFour = treeModel.Get(treeModel.MangleNum(lib, fourNum)) isnt: #())
		treeModel.Rename(recFour, 'Three')
		Assert(recFour = svcTable.Get('Three') isnt: false)
		Assert(recFour.lib_before_text is: 'three text')

		Assert(svc.Local_changes(lib).Map({ it.type $ ' ' $ it.name })
			is: #('  Four', '  Three'))

		// rename Four to xFour, Three to Four, xFour to Three
		Assert(recThree = treeModel.Get(treeModel.MangleNum(lib, threeNum)) isnt: #())
		treeModel.Rename(recThree, 'xFour')
		Assert(recThree = svcTable.Get('xFour') isnt: false)
		Assert(recThree.lib_committed is: '')

		Assert(recFour = treeModel.Get(treeModel.MangleNum(lib, fourNum)) isnt: #())
		treeModel.Rename(recFour, 'Four')
		Assert(recFour = svcTable.Get('Four') isnt: false)
		Assert(recFour.lib_before_text is: '')
		Assert(recFour.lib_committed isnt: '')

		Assert(recThree = treeModel.Get(treeModel.MangleNum(lib, threeNum)) isnt: #())
		treeModel.Rename(recThree, 'Three')
		Assert(recThree = svcTable.Get('Three') isnt: false)
		Assert(recThree.lib_before_text is: '')
		Assert(recThree.lib_committed isnt: '')
		}

	Test_SetModified()
		{
		svcTable = .setupLib([
			[group: 0, parent: 0, name: #Folder1, num: 1],
			[group: 0, parent: 0, name: #Folder2, num: 2],
			[name: #Rec, text: 'text', lib_before_text: 'text',
				lib_modified: Date(), lib_committed: Date(), parent: 1]
			])
		table = svcTable.Table()

		treeModel = .treeModel([table])
		num = treeModel.MangleNum(table, Query1(table, name: #Rec).num)
		rec = treeModel.Get(num)
		folder1Rec = Query1(table, name: #Folder1, group: 0)
		folder2Rec = Query1(table, name: #Folder2, group: 0)

		// record is modified
		rec = treeModel.Get(num)
		rec.text = 'text modified'
		treeModel.Update(rec)
		Assert(saved = svcTable.Get(rec.name) isnt: false)
		Assert(saved.lib_modified isnt: '')

		// record has been modified back to the original
		rec = treeModel.Get(num)
		rec.text = 'text'
		treeModel.Update(rec)
		Assert(saved = svcTable.Get(rec.name) isnt: false)
		Assert(saved.lib_modified is: '')

		// record was made invalid, prompting CodeState to handle it
		rec = treeModel.Get(num)
		rec.lib_invalid_text = 'invalid'
		treeModel.Update(rec)
		Assert(saved = treeModel.Get(num) isnt: false)
		Assert(saved.lib_modified isnt: '')

		// record is restored to original state
		rec = treeModel.Get(num)
		rec.text = 'text'
		rec.lib_invalid_text = ''
		treeModel.Update(rec)
		Assert(saved = treeModel.Get(num) isnt: false)
		Assert(saved.lib_modified is: '')

		// record is moved to a new folder
		rec = treeModel.Get(num)
		treeModel.Move(rec, newParent = treeModel.MangleNum(table, folder2Rec.num))
		Assert(saved = treeModel.Get(num) isnt: false)
		Assert(saved.lib_modified isnt: '')
		Assert(saved.parent is: newParent)
		Assert(saved.lib_before_path is: folder1Rec.name)

		// record is moved back to original folder, text is different
		num = treeModel.MangleNum(table, Query1(table, name: #Rec).num)
		rec = treeModel.Get(num)
		rec.text = 'text changed'
		treeModel.Update(rec)
		treeModel.Move(treeModel.Get(num),
			newParent = treeModel.MangleNum(table, folder1Rec.num))
		Assert(saved = treeModel.Get(num) isnt: false)
		Assert(saved.lib_modified isnt: '')
		Assert(saved.parent is: newParent)
		Assert(saved.lib_before_path is: '')

		// record is restored to base
		rec.text = 'text'
		treeModel.Update(rec)
		Assert(saved = treeModel.Get(num) isnt: false)
		Assert(saved.lib_modified is: '')
		}

	setupLib(recs)
		{
		recs.Map!({ it.lib_committed = Date(); it })
		svcTable = .SvcTable(recs.table = .MakeLibrary())
		.MakeLibraryRecord(@recs)
		recs.Delete('table').Each({ it.path = svcTable.GetPath(it) })
		return svcTable
		}

	Test_setItemText()
		{
		treeModel = .treeModel([.MakeLibrary()])
		m = treeModel.LibTreeModel_setItemText

		item = []
		m(item, 0, origText?: false)
		Assert(item.lib_before_text is: '')

		item.text = 'text'
		m(item, 0, origText?: false)
		Assert(item.lib_before_text is: '')

		item.lib_committed = Date()
		m(item, 0, origText?: false)
		Assert(item.lib_before_text is: 'text')

		item.text = 'text2'
		m(item, 0, origText?: false)
		Assert(item.lib_before_text is: 'text')

		item.lib_before_text = ''
		item.lib_modified = Date()
		m(item, 0, origText?: false)
		Assert(item.lib_before_text is: '')

		item.lib_invalid_text = 'invalid code'
		item.lib_modified = ''
		m(item, 0, origText?: false)
		Assert(item.text is: 'invalid code')
		Assert(item.lib_before_text is: 'text2')

		item.text = 'text2'
		item.lib_before_text = ''
		m(item, 0, origText?:)
		Assert(item.text is: 'text2')
		Assert(item.lib_before_text is: 'text2')
		}

	Test_move()
		{
		svcLib1 = .setupLib(lib1Recs = [
			[group: 0, parent: 0, num: 1, name: #Folder1],
			[group: 0, parent: 0, num: 2, name: #Folder2],
				[parent: 2, num: 4, name: #Rec1],
				[group: 2, parent: 2, num: 3, name: #Folder3],
					[parent: 3, num: 5, name: #Rec2],
			[parent: 0, num: 6, name: #RecX],
			])
		svcLib2 = .setupLib(lib2Recs = [
			[group: 0, parent: 0, num: 4, name: #Folder3],
				[parent: 4, num: 5, name: #Rec3],
			[parent: 0, num: 7, name: #RecX]
			])

		// -- Moving folders --
		.testMoves(#Folder1, svcLib1, svcLib1, #Folder2)
		.testMoves(#Folder1, svcLib1, svcLib1, #Folder3)
		.testMoves(#Folder1, svcLib1, svcLib2, #Folder3)
		.testMoves(#Folder1, svcLib1, svcLib2, #root)

		.testMoves(#Folder2, svcLib1, svcLib1, #Folder1)
		.testMoves(#Folder2, svcLib1, svcLib2, #Folder3)
		.testMoves(#Folder2, svcLib1, svcLib2, #root)

		.testMoves(#Folder3, svcLib1, svcLib1, #Folder1)
		.testMoves(#Folder3, svcLib1, svcLib1, #root)
		.testMoves(#Folder3, svcLib1, svcLib2, #Folder3)

		.testMoves(#Folder3, svcLib2, svcLib1, #root)
		.testMoves(#Folder3, svcLib2, svcLib1, #Folder1)
		.testMoves(#Folder3, svcLib2, svcLib1, #Folder3)

		// -- Moving records --
		.testMoves(#Rec1, svcLib1, svcLib1, #root)
		.testMoves(#Rec1, svcLib1, svcLib1, #Folder1)
		.testMoves(#Rec1, svcLib1, svcLib1, #Folder3)
		.testMoves(#Rec1, svcLib1, svcLib2, #root)
		.testMoves(#Rec1, svcLib1, svcLib2, #Folder3)

		.testMoves(#Rec2, svcLib1, svcLib1, #root)
		.testMoves(#Rec2, svcLib1, svcLib1, #Folder1)
		.testMoves(#Rec2, svcLib1, svcLib1, #Folder2)
		.testMoves(#Rec2, svcLib1, svcLib1, #Folder3)
		.testMoves(#Rec2, svcLib1, svcLib2, #root)

		// -- Verify tables states, should match the original state at this point --
		.assertTableState(svcLib1, lib1Recs)
		.assertTableState(svcLib2, lib2Recs)

		// -- Special cases --
		// Ensure unique record names when moving Folders between tables
		// When moving Folder3 to the root of the other table it will be renamed to
		// stay unique
		.testMoves(#Folder3, svcLib1, svcLib2, #root)
		Assert(Query1(svcLib1.Table(), name: 'Folder3_Copy1') isnt: false)

		// Now that Folder3 in svcLib1 is Folder3_Copy1 (after moving back), Folder3 from
		// svclib2 does not require a rename
		.testMoves(#Folder3, svcLib2, svcLib1, #Folder2)
		Assert(Query1(svcLib2.Table(), name: 'Folder3') isnt: false)

		// Record will be automatically renamed to ensure it does not conflict with the
		// other RecX
		.testMoves(#RecX, svcLib2, svcLib1, #root)
		Assert(Query1(svcLib2.Table(), name: 'RecX_Copy1') isnt: false)
		}

	testMoves(moveName, fromTable, toTable, moveToName)
		{
		// Carry out the move / restore and assert the results
		move = Query1(fromTable.Table(), name: moveName)
		move.group = move.group > -1
		moveFrom = move.parent
		moveTo = moveToName isnt 'root'
			? Query1(toTable.Table(), name: moveToName).num
			: 0
		.assertMove(move, moveFrom, moveTo, fromTable, toTable)
		.assertMove(move, moveTo, moveFrom, toTable, fromTable)
		}

	assertMove(move, moveFrom, moveTo, fromTable, toTable)
		{
		origNum = move.num
		childrenBefore = QueryAll(fromTable.Table(), parent: move.num)
		// Move folder
		Assert(move.parent is: moveFrom)
		Transaction(update:)
			{ |t|
			LibTreeModel.
				LibTreeModel_move(move, moveTo, fromTable, toTable, t)
			}
		Assert(move.parent is: moveTo)
		// If not a library move, the record nums should NOT change
		if fromTable.Table() is toTable.Table()
			Assert(origNum is: move.num)
		// Verify that ALL children records moved as well
		childrenAfter = QueryAll(toTable.Table(), parent: move.num)
		for child in childrenBefore
			Assert(childrenAfter.Any?({ it.name is child.name}),
				msg: 'Failed to move: ' $ child.name)
		}

	assertTableState(svcTable, recs)
		{
		recs.Each()
			{
			rec = Query1(svcTable.Table(), name: it.name)
			Assert(svcTable.GetPath(rec) is: it.path)
			}
		}

	Test_renameFolders()
		{
		svcTable = .setupLib(recs = [
			folder1 = [group: 0, parent: 0, num: 1, name: #Folder1],
				[parent: 1, num: 4, name: #Rec1],
				folder2 = [group: 1, parent: 1, num: 2, name: #Folder2],
					[parent: 2, num: 5, name: #Rec2],
					folder3 = [group: 2, parent: 2, num: 3, name: #Folder3],
						[parent: 3, num: 6, name: #Rec3]
			])

		table = svcTable.Table()
		Assert(Query1(table, num: 4).lib_before_path is: '')
		Assert(Query1(table, num: 5).lib_before_path is: '')
		Assert(Query1(table, num: 6).lib_before_path is: '')

		.assertRename(folder1, svcTable, recs, #Folder4)
		Assert(Query1(table, num: 4).lib_before_path is: 'Folder1')
		Assert(Query1(table, num: 5).lib_before_path is: 'Folder1/Folder2')
		Assert(Query1(table, num: 6).lib_before_path is: 'Folder1/Folder2/Folder3')

		.assertRename(folder2, svcTable, recs, #Folder5)
		Assert(Query1(table, num: 4).lib_before_path is: 'Folder1')
		Assert(Query1(table, num: 5).lib_before_path is: 'Folder1/Folder2')
		Assert(Query1(table, num: 6).lib_before_path is: 'Folder1/Folder2/Folder3')

		.assertRename(folder3, svcTable, recs, #Folder6)
		Assert(Query1(table, num: 4).lib_before_path is: 'Folder1')
		Assert(Query1(table, num: 5).lib_before_path is: 'Folder1/Folder2')
		Assert(Query1(table, num: 6).lib_before_path is: 'Folder1/Folder2/Folder3')

		.assertRename(folder1, svcTable, recs, #Folder1, renamedBack?:)
		Assert(Query1(table, num: 4).lib_before_path is: '')
		Assert(Query1(table, num: 5).lib_before_path is: 'Folder1/Folder2')
		Assert(Query1(table, num: 6).lib_before_path is: 'Folder1/Folder2/Folder3')

		.assertRename(folder2, svcTable, recs, #Folder2, renamedBack?:)
		Assert(Query1(table, num: 4).lib_before_path is: '')
		Assert(Query1(table, num: 5).lib_before_path is: '')
		Assert(Query1(table, num: 6).lib_before_path is: 'Folder1/Folder2/Folder3')

		.assertRename(folder3, svcTable, recs, #Folder3, renamedBack?:)
		Assert(Query1(table, num: 4).lib_before_path is: '')
		Assert(Query1(table, num: 5).lib_before_path is: '')
		Assert(Query1(table, num: 6).lib_before_path is: '')
		}

	// Validate all children (recs object is ordered in folder nesting order)
	assertRename(folderRec, svcTable, recs, newName, renamedBack? = false)
		{
		table = svcTable.Table()
		LibTreeModel.LibTreeModel_renameFolder(folderRec, table, newName)
		// Ensure rename worked
		Assert(Query1(table, num: folderRec.num).name is: newName)

		for rec in recs[recs.FindIf({ it.name is folderRec.name }) + 1 ..]
			{
			lookup = Query1(table, num: rec.num)
			lookup.path = svcTable.GetPath(lookup)
			// If the record is not a folder record, verify lib_before_path is set
			if not rec.Member?(#group)
				Assert(lookup.lib_before_path
					is: renamedBack? and rec.parent is folderRec.num ? '' : rec.path)
			Assert(rec.num is: lookup.num)
			Assert(rec.parent is: lookup.parent)
			}
		}

	Test_Modified?()
		{
		.SpyOn(SvcSettings.Set?).Return(false, true, true, true, true)
		m = LibTreeModel.Modified?

		// No SvcSettings, record is never seen as modified
		rec = [lib_committed: '', lib_modified: '', group: false]
		Assert(m(rec) is: false)

		// SvcSettings are filled, record is new
		Assert(m(rec))

		// New Record is committed
		rec.lib_committed = Date()
		Assert(m(rec) is: false)

		// Record is modified
		rec.lib_modified = Date()
		Assert(m(rec))

		// Record is a folder
		rec.group = true
		Assert(m(rec) is: false)

		// Modifications are removed
		rec.group = false
		rec.lib_modified = ''
		Assert(m(rec) is: false)

		// Simulate root folder, never seen as modified
		Assert(m([]) is: false)
		}

	Test_Valid?()
		{
		m = LibTreeModel.Valid?

		Assert(m([]))
		Assert(m([lib_invalid_text: '']))
		Assert(m([lib_invalid_text: 'text']) is: false)
		}

	Test_Synced?()
		{
		m = LibTreeModel.Synced?

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

		rec.text = 'function () { /*comment*/ }'
		Assert(m(rec, savedRec) is: false)

		savedRec.text = 'function () { /*comment*/ }'
		Assert(m(rec, savedRec))

		savedRec.text = 'function () { /*Comment*/ }'
		Assert(m(rec, savedRec) is: false)
		}

	// NOTE: We do not test deleteLibrary as it produces issues with various Library
	// components. These issues render work-copies borderline unusable
	Test_DeleteItem()
		{
		.SpyOn(SvcTable.Publish).Return('')
		svcTable = .setupLib([
			[group: 0, parent: 0, num: 1, name: #Folder1],
			[parent: 0, num: 9, name: #Rec0],
				[parent: 1, num: 4, name: #Rec1],
				[group: 1, parent: 1, num: 2, name: #Folder2],
					[parent: 2, num: 5, name: #Rec2],
					[group: 2, parent: 2, num: 3, name: #Folder3_1],
						[parent: 3, num: 6, name: #Rec3_1],
					[group: 2, parent: 2, num: 7, name: #Folder3_2],
						[parent: 7, num: 8, name: #Rec3_2]
			])
		table = svcTable.Table()
		treeModel = .treeModel([table])

		// Delete a record, should be "staged" for deletion
		treeModel.DeleteItem(treeModel.MangleNum(table, 8), #Rec3_2, false)
		Assert(svcTable.Get(#Rec3_2) is: false)
		Assert(rec = svcTable.Get(#Rec3_2, deleted:) isnt: false)
		Assert(rec.lib_before_path is: 'Folder1/Folder2/Folder3_2')

		// Delete a empty folder, folder should be removed
		treeModel.DeleteItem(treeModel.MangleNum(table, 7), #Folder3_2, true)
		Assert(Query1(table, num: 7) is: false)
		Assert(QueryEmpty?(table, parent: 7, group: -1), msg: 'delete empty folder')

		// Delete a folder with a record. Record should be deleted as well
		treeModel.DeleteItem(treeModel.MangleNum(table, 3), #Folder3_1, true)
		Assert(Query1(table, num: 3) is: false)
		Assert(QueryEmpty?(table, parent: 3, group: -1), msg: 'delete foler with 1 rec')
		Assert(svcTable.Get(#Rec3_1) is: false)
		Assert(rec = svcTable.Get(#Rec3_1, deleted:) isnt: false)
		Assert(rec.lib_before_path is: 'Folder1/Folder2/Folder3_1')

		// Delete a folder with nested folders / records.
		// All nested records should be removed
		treeModel.DeleteItem(treeModel.MangleNum(table, 1), #Folder1, true)
		Assert(Query1(table, num: 1) is: false)
		Assert(QueryEmpty?(table, parent: 1, group: -1), msg: 'delete with sub folder')
		Assert(svcTable.Get(#Rec1) is: false)
		Assert(rec = svcTable.Get(#Rec1, deleted:) isnt: false)
		Assert(rec.lib_before_path is: 'Folder1')
		// Folder2 should be deleted along with its children
		Assert(Query1(table, num: 2) is: false)
		Assert(QueryEmpty?(table, parent: 2, group: -1), msg: 'folder 2 with children')
		Assert(svcTable.Get(#Rec2) is: false)
		Assert(rec = svcTable.Get(#Rec2, deleted:) isnt: false)
		Assert(rec.lib_before_path is: 'Folder1/Folder2')

		// Delete a record, should be "staged" for deletion (root level)
		treeModel.DeleteItem(treeModel.MangleNum(table, 9), #Rec0, false)
		Assert(svcTable.Get(#Rec0) is: false)
		Assert(rec = svcTable.Get(#Rec0, deleted:) isnt: false)
		Assert(rec.lib_before_path is: '<root>')
		}

	Test_TreeSort()
		{
		rec1 = [name: 'A', group:]
		rec2 = [name: 'B', group:]
		mock = Mock(LibTreeModel)
		mock.When.Get(rec1).Return(rec1)
		mock.When.Get(rec2).Return(rec2)
		mock.When.TreeSort([anyArgs:]).CallThrough()

		// rec1 comes first as both are folders, and rec1 is alphabetically first
		Assert(mock.TreeSort(rec1, rec2) is: -1)

		// rec2 comes first as its a folder, and rec1 is not
		rec1.group = false
		Assert(mock.TreeSort(rec1, rec2) is: 1)

		// rec1 comes first as both are folders, and rec1 is alphabetically first
		rec2.group = false
		Assert(mock.TreeSort(rec1, rec2) is: -1)

		// rec2 comes first as both are folders, and rec2 is ascii first
		rec1.name = 'a'
		Assert(mock.TreeSort(rec1, rec2) is: 1)
		}
	}
