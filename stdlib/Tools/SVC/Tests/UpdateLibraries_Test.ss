// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		classMock = Mock(UpdateLibraries)
		classMock.When.closeSvcSocketClient().Do({ })
		classMock.When.log([anyArgs:]).Do({ })
		classMock.When.updateLibraries([anyArgs:]).CallThrough()
		classMock.When.getChanges().CallThrough()

		// Svc error handling
		// Error reported during Svc init, error is thrown and returned
		errorPrefix = 'ERROR: UpdateLibraries: '
		err = 'Unexpected Issue'
		classMock.When.svc().Return(err, svcMock = Mock())
		Assert(classMock.updateLibraries(#(), false) is: errorPrefix $ err)
		// Error reported by CheckSvcStatus, error is returned
		err = 'ERR SocketClient is not connected'
		svcMock.When.CheckSvcStatus().Return(err, '')
		Assert(classMock.updateLibraries(#(), false)
			is: errorPrefix $ 'Svc status: ' $ err)

		// No errors, tables object is empty so there are no tables to process
		Assert(classMock.updateLibraries(#(), false) is: 0)

		// Tables object only contains a non-existent table
		tables = [fakeLib = .TempTableName()]
		Assert(classMock.updateLibraries(tables, false) is: 0)
		classMock.Verify.log('WARNING: ' $ fakeLib $ ' does not exist', '')

		// Tables object contains a table which exists, changes and update (7 records)
		tables.Add(newLib = .MakeLibrary([]))
		svcMock.When.GetChanges(newLib).Return([
			conflicts: 		[],
			local_changes: 	[],
			master_changes: [newLib] // Used to control UpdateLibrary return value
			])
		svcMock.When.UpdateLibrary([newLib]).Return(7)
		Assert(classMock.updateLibraries(tables, false) is: 7)
		classMock.Verify.Times(2).log([anyArgs:])
		classMock.Verify.Times(2).log('WARNING: ' $ fakeLib $ ' does not exist', '')

		// Tables object contains a table with conflicts
		tables.Add(conflictLib = .MakeLibrary([]))
		svcMock.When.GetChanges(conflictLib).Return([
			conflicts: 		[[name: #Con1], [name: #Con2]],
			local_changes: 	[[name: #Chg1], [name: #Chg2]],
			master_changes: [conflictLib]
			])
		Assert(classMock.updateLibraries(tables, false)
			is: 'FAILURES: Version Control Conflicts in ' $ conflictLib $ '\n')
		classMock.Verify.Times(4).log([anyArgs:])
		classMock.Verify.Times(3).log('WARNING: ' $ fakeLib $ ' does not exist', '')
		classMock.Verify.log(errorPrefix $ conflictLib $ ' has conflicts', #(Con1, Con2))

		// Tables object contains a table with conflicts
		tables.Add(localChangesLib = .MakeLibrary([]))
		svcMock.When.GetChanges(localChangesLib).Return([
			conflicts: 		[],
			local_changes: 	[[name: #Chg3], [name: #Chg4]],
			master_changes: [localChangesLib]
			])
		svcMock.When.UpdateLibrary([localChangesLib]).Return(3)
		Assert(classMock.updateLibraries(tables, false)
			is: 'FAILURES: Version Control Conflicts in ' $ conflictLib $ '\n')
		classMock.Verify.Times(7).log([anyArgs:])
		classMock.Verify.Times(4).log('WARNING: ' $ fakeLib $ ' does not exist', '')
		classMock.Verify.Times(2).
			log(errorPrefix $ conflictLib $ ' has conflicts', #(Con1, Con2))
		classMock.Verify.
			log('WARNING: ' $ localChangesLib $ ' has local changes', #(Chg3, Chg4))

		// Removing conflictLib from tables. Should return the number of changes processed
		tables.Remove(conflictLib)
		Assert(classMock.updateLibraries(tables, false) is: 10)
		classMock.Verify.Times(9).log([anyArgs:])
		classMock.Verify.Times(5).log('WARNING: ' $ fakeLib $ ' does not exist', '')
		classMock.Verify.Times(2).
			log(errorPrefix $ conflictLib $ ' has conflicts', #(Con1, Con2))
		classMock.Verify.Times(2).
			log('WARNING: ' $ localChangesLib $ ' has local changes', #(Chg3, Chg4))

		// Adding a book table to tables and testing preGetFn
		tables.Add(book = .MakeBook())
			svcMock.When.GetChanges(book).Return([
			conflicts: 		[],
			local_changes: 	[],
			master_changes: [book]
			])
		preGetFn = {
			|masterChanges, unused|
			Assert(masterChanges[0] isnt: book)
			}
		svcMock.When.UpdateLibrary([book]).Return(9)
		Assert(classMock.updateLibraries(tables, preGetFn) is: 19)
		classMock.Verify.Times(11).log([anyArgs:])
		classMock.Verify.Times(6).log('WARNING: ' $ fakeLib $ ' does not exist', '')
		classMock.Verify.Times(2).
			log(errorPrefix $ conflictLib $ ' has conflicts', #(Con1, Con2))
		classMock.Verify.Times(3).
			log('WARNING: ' $ localChangesLib $ ' has local changes', #(Chg3, Chg4))
		}
	}
