// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	svcUser: SvcUsers
		{
		New() { }

		SetTables(userTable, passwordTable)
			{
			.SvcUsers_svcPasswordTable = passwordTable
			.SvcUsers_svcUsersTable = userTable
			}

		key: false
		SvcUsers_send(sc /*unused*/, args)
			{
			cmd = args.Extract(0)
			args = [.key].Merge(args)
			result = false
			switch (cmd)
				{
				case 'LOGIN':
					result = this.LoginRequest(@args)
				case 'ADDUSER':
					result = this.AddUserRequest(@args)
				case 'CHANGEPASSWORD':
					result = this.ChangePasswordRequest(@args)
				case 'DELETEUSER':
					result = this.DeleteUserRequest(@args)
				case 'NONCE':
					result = .key = SvcServer.SvcServer_buildKey()
				}
			return result
			}

		SvcUsers_ensure() { }

		SvcUsers_socketClient()
			{ return [] }
		}

	Test_main()
		{
		cl = new .svcUser

		// ----- Server Validation -----
		// Password table does not exist
		cl.SetTables('fake_table', 'fake_table')
		Assert(cl.AddUser('wrongPassword', 'user', 'password')
			is: 'Server is not secure')
		// Password table exists, no password is set
		cl.SetTables(
			userTable = .MakeTable('(svcuser_id, svcuser_passhash) key (svcuser_id)'),
			passTable = .MakeTable('(svc_password) key (svc_password)'))
		Assert(cl.AddUser('wrongPassword', 'user', 'password')
			is: 'Server is not secure')
		// Password table is set / ready, wrong server password
		QueryOutput(passTable, [svc_password: PassHash('',  serverPassword = 'test')])
		Assert(cl.AddUser('wrongPassword', 'user', 'password')
			is: 'Invalid Server Password')

		// ----- Add New User -----
		// User is added successfully
		Assert(cl.AddUser(serverPassword, 'user', 'password'))
		Assert(Query1(userTable, svcuser_id: 'user') isnt: false)
		// Attempt to add the user again, validation stops it
		Assert(cl.AddUser(serverPassword, 'user', 'password') is: 'User already exists')

		// ----- Login, Post User Creation -----
		// User attempts login, wrong user name
		Assert(cl.TestLogin('wrongUser', 'wrongPassword') is: false)
		// User attempts login, wrong password
		Assert(cl.TestLogin('user', 'wrongPassword') is: false)
		// User logs in successfully
		Assert(cl.TestLogin('user', 'password'))

		// ----- Password Change -----
		// Attempt to change password, wrong server password
		res = cl.ChangePassword('wrongPassword', 'user2', 'password', 'newPassword')
		Assert(res is: 'Invalid Server Password')
		// Attempt to change password, non-existent user
		res = cl.ChangePassword(serverPassword, 'user2', 'password', 'newPassword')
		Assert(res is: false)
		// Attempt to change password, wrong old password is used
		res = cl.ChangePassword(serverPassword, 'user', 'wrongOldPassword', 'newPassword')
		Assert(res is: 'Old password does not match the saved password')
		// Attempt to change password, correct old password is used
		res = cl.ChangePassword(serverPassword, 'user', 'password', 'newPassword')
		Assert(res)

		// ----- Login, Post Password Change -----
		// User attempts login, old password
		Assert(cl.TestLogin('user', 'password') is: false)
		// User logs in with new password
		Assert(cl.TestLogin('user', 'newPassword'))

		// ----- Delete User -----
		// Attempt to delete user, wrong server password
		Assert(cl.DeleteUser('wrongPassword', 'user') is: 'Invalid Server Password')
		// Attempt to delete user, correct server password
		Assert(cl.DeleteUser(serverPassword, 'user'))
		Assert(Query1(userTable, svcuser_id: 'user') is: false)
		// Attempt to delete non-existent user, correct server password
		Assert(cl.DeleteUser(serverPassword, 'user') is: false)
		}
	}