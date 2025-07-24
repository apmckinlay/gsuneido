// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.user = Suneido.User
		Suneido.User = 'default'
		.created = false is Query1('tables', table: UserSettings.Table)
		}

	key: 'BUSTestKey'
	Test_main()
		{
		Assert(UserSettings.Get(.key) is: false)
		Assert(UserSettings.Get(.key, 'xxx') is: 'xxx')

		UserSettings.Put(.key, 'Hello')
		Assert(UserSettings.Get(.key) is: 'Hello')

		UserSettings.Put(.key, 'World')
		Assert(UserSettings.Get(.key) is: 'World')

		UserSettings.Put(.key, 'UserTest', 'testUser')
		Assert(UserSettings.Get(.key, user: 'testUser') is: 'UserTest')

		UserSettings.Put(.key, 'UserTest', 'testUser')
		Assert(UserSettings.Get(.key, user: 'diffUser') is: false)
		}

	Teardown()
		{
		UserSettings.RemoveAllUsers(.key)
		Suneido.User = .user
		if .created
			Database('destroy ' $ UserSettings.Table)
		}
	}