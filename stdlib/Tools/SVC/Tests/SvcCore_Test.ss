// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
SvcTests
	{
	Test_Exists()
		{
		Assert(SvcCore.Exists?(lib = .MakeLibrary()) is: false)
		.CommitAdd(.Svc(), .SvcTable(lib), .TempName(), 'text', 'add')
		Assert(SvcCore.Exists?(lib))
		}

	Test_OnlyDeletedChangesBetween()
		{
		m = SvcCore.OnlyDeletedChangesBetween?
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		added1 = .CommitAdd(svc, svcTable, rec1 = .TempName(), 'text', 'add')
		deleted1 = .CommitDelete(svc, svcTable, rec1, 'delete')
		added2 = .CommitAdd(svc, svcTable, rec2 = .TempName(), 'text', 'add')

		Assert(m(lib, added1, Timestamp()) is: false)
		Assert(m(lib, added1, deleted1))
		Assert(m(lib, deleted1, added2) is: false)

		deleted2 = .CommitDelete(svc, svcTable, rec2, 'delete')
		Assert(m(lib, added2, deleted2))
		}

	Test_SearchForRename()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())

		// Record is output
		lib_before_hash = Svc.Hash(text = 'Original Text: ' $ Display(Timestamp()))
		.CommitAdd(svc, svcTable, 'OrigRec', text, 'add')
		// Record is renamed
		.CommitDelete(svc, svcTable, 'OrigRec', 'Rename OrigRec to Rename1')
		.CommitAdd(svc, svcTable, 'Rename1', text, 'Rename OrigRec to Rename1',
			:lib_before_hash)
		Assert(SvcCore.SearchForRename(lib, 'Rename1') is: #('OrigRec'))
		Assert(SvcCore.SearchForRename(lib, 'OrigRec').Empty?())

		// Record is renamed again, Order of delete / add should NOT matter
		.CommitAdd(svc, svcTable, 'Rename2', text, 'Rename OrigRec to Rename1',
			:lib_before_hash)
		.CommitDelete(svc, svcTable, 'Rename1', 'Rename Rename1 to Rename2')
		Assert(SvcCore.SearchForRename(lib, 'Rename2') is: #('OrigRec', 'Rename1'))
		}

	Test_GetBeforeMethods()
		{
		svc = .Svc()
		svcTable = .SvcTable(lib = .MakeLibrary())
		.CommitAdd(svc, svcTable, 'Rec', '', 'add')
		bef = Timestamp()
		for i in ..10
			{
			.CommitTextChange(svc, svcTable, 'Rec', Display(i), text = 'modify - ' $ i)
			Assert(SvcCore.GetBefore(lib, 'Rec', Timestamp()).comment is: text)
			}
		Assert(SvcCore.GetBefore(lib, 'Rec', bef).comment is: 'add')
		Assert(SvcCore.GetBefore(lib, 'Rec', bef).comment is: 'add')

		Assert(result = SvcCore.Get10Before(lib, 'Rec', Timestamp()) isSize: 10)
		Assert(result.Any?({ it.comment is 'add' }) is: false)

		Assert(result = SvcCore.Get10Before(lib, 'Rec', bef) isSize: 1)
		Assert(result[0].comment is: 'add')
		}

	Test_ensureMaster_masterType()
		{
		ensure = SvcCore.EnsureMaster
		masterType = SvcCore.MasterType
		table = .TempTableName()
		type = 'notStandard'
		errMsg = 'ERR ensuring master table: ' $ table $ '_master, invalid type: ' $ type
		Assert(ensure(table $ '_master', type) is: errMsg)
		Assert(masterType(table) is: false)

		Assert(ensure((table = .TempTableName()) $ '_master', 'lib'))
		Assert(masterType(table) is: 'lib')

		Assert(ensure((table = .TempTableName()) $ '_master', 'book'))
		Assert(masterType(table) is: 'book')
		}
	}