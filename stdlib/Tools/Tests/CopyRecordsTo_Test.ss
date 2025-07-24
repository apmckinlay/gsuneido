// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		name1 = .TempName()
		name2 = .TempName()
		rec = [name: name1, text: 'class {}']
		srcLib = .MakeLibrary(rec)
		midLib = .MakeLibrary(rec)
		dstLib = .MakeLibrary()
		mock = Mock(CopyRecordsTo)
		mock.When.libraries().Return([srcLib, midLib, dstLib])
		mock.When.checkIfOverloaded([anyArgs:]).CallThrough()
		mock.When.createDstFolder([anyArgs:]).CallThrough()
		mock.When.copyRecords([anyArgs:]).CallThrough()
		mock.When.nextNum([anyArgs:]).CallThrough()

		fn = { |srcLib, srcNames, dstLib, dstFolder, overwrite?|
			mock.Eval(CopyRecordsTo.CallClass, srcLib, srcNames, dstLib, dstFolder,
				function(@unused){}, overwrite?)
			}
		// no record to copy
		Assert(fn(srcLib, [], dstLib, '', false) is: 0)

		// record is overloaded in midLib
		Assert(fn(srcLib, [name1], dstLib, '', false)
			is: name1 $ ' is overloaded in ' $ midLib)

		mock.When.libraries().Return([srcLib, dstLib])
		// record doesn't exist in srcLib
		Assert(fn(srcLib, [name1, name2], dstLib, '', false)
			is: name2 $ ' not found in ' $ srcLib)
		Assert(Query1(dstLib, name: name1, group: -1) is: false)

		// move to a different folder
		Assert(fn(srcLib, [name1], srcLib, 'test2', false) is: 1)
		newParent = Query1(srcLib, name: 'test2').num
		.check(srcLib, name1, newParent) // Note: parent was changed to the new folder

		// copy to a different library
		Assert(fn(srcLib, [name1], dstLib, 'test', false) is: 1)
		parent = Query1(dstLib, name: 'test').num
		.check(dstLib, name1, parent)

		// overwrite records in original place
		Assert(fn(srcLib, [name1], dstLib, 'test2', true) is: 1)
		.check(dstLib, name1, parent)

		// move records to specified folder
		Assert(fn(srcLib, [name1], dstLib, 'test2', false) is: 1)
		newParent = Query1(dstLib, name: 'test2').num
		.check(dstLib, name1, newParent)

		// delect and copy to default folder
		QueryDelete(dstLib, Query1(dstLib, name: name1, group: -1))
		// parent should be 0 when dstFolder is ''
		Assert(fn(srcLib, [name1], dstLib, '', false) is: 1)
		.check(dstLib, name1, 0)
		}

	check(lib, name, parent = false)
		{
		rec = Query1(lib, :name, group: -1)
		Assert(rec isnt: false)
		if parent isnt false
			Assert(rec.parent is: parent)
		}
	}
