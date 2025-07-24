// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_updateSettings()
		{
		mock = Mock(SvcSettings)
		mock.When.updateSettings([anyArgs:]).CallThrough()
		mock.When.columns().Return(#(svc_server, svc_user, svc_local?))
		mock.When.update([anyArgs:]).Do({ })

		mock.updateSettings(false, [])
		mock.Verify.columns()
		mock.Verify.Never().update([anyArgs:])

		mock.updateSettings([], [])
		mock.Verify.Times(2).columns()
		mock.Verify.Never().update([anyArgs:])

		mock.updateSettings([svc_user: 'user'], [svc_user: 'user'])
		mock.Verify.Never().update([anyArgs:])

		mock.updateSettings([], [svc_user: 'user'])
		mock.Verify.update([svc_user: 'user'])

		mock.updateSettings([svc_local?:], [svc_user: 'user'])
		mock.Verify.Times(2).update([svc_user: 'user'])
		}

	Test_UpdateCredentials()
		{
		.SpyOn(DeleteFile).Return(0)
		.SpyOn(PutFile).Return(0)

		mock = Mock(SvcSettings)
		mock.When.UpdateCredentials([anyArgs:]).CallThrough()
		Assert(mock.UpdateCredentials('', '') is: [])

		mock.When.decrypt([anyArgs:]).Return([userId: #test, passhash: #hash])
		Assert(mock.UpdateCredentials(#test, #hash) is: [userId: #test, passhash: #hash])
		mock.Verify.Never().encodeDecode([anyArgs:])

		Assert(mock.UpdateCredentials(#test1, #hash)
			is: [userId: #test1, passhash: PassHash(#test1, #hash)])
		mock.Verify.encodeDecode([anyArgs:])

		Assert(mock.UpdateCredentials(#test, #hash1)
			is: [userId: #test, passhash: PassHash(#test, #hash1)])
		mock.Verify.Times(2).encodeDecode([anyArgs:])
		}

	Test_Set?()
		{
		mock = Mock(SvcSettings)
		mock.Table = .TempTableName()
		mock.When.Set?().CallThrough()
		mock.When.Ensure().CallThrough()
		teardown = {
			if TableExists?(mock.Table)
				Database('drop ' $ mock.Table)
			}
		.AddTeardown(teardown)

		// No table
		Assert(TableExists?(mock.Table) is: false)
		Assert(mock.Set?() is: false)

		// Empty table
		mock.Ensure()
		Assert(TableExists?(mock.Table))
		Assert(mock.Set?() is: false)

		// Empty record
		QueryOutput(mock.Table, [])
		Assert(mock.Set?() is: false)

		// Server is set
		QueryDo('update ' $ mock.Table $ ' set svc_server = "set"')
		Assert(mock.Set?())

		// Local is set
		QueryDo('update ' $ mock.Table $ ' set svc_server = "", svc_local? = true')
		Assert(mock.Set?())

		// Neither is set (svc_local? is false instead of "")
		QueryDo('update ' $ mock.Table $ ' set svc_server = "", svc_local? = false')
		Assert(mock.Set?() is: false)
		}
	}
