// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.created = QueryEmpty?('tables', table: UserSettings.Table)
		UserSettingsCached.Reset()
		.key = .TempName()
		.key2 = .TempName()
		.key3 = .TempName()
		.user1 = .TempName()
		.user2 = .TempName()
		}

	Test_Cache()
		{

		// user1 gets default value
		result = UserSettingsCached().Get(.key, 'defaultValue', .user1)
		cache = Suneido.GetDefault(
			'Memoize_UserSettingsCached.UserSettingsCached_cache', false)
		Assert(cache isnt: false)
		Assert(cache.GetMissRate() is: 1/1) // 1 miss / 1 get
		Assert(result is: 'defaultValue')

		// user1 gets default value again. This time it is cached
		result = UserSettingsCached().Get(.key, 'defaultValue', .user1)
		Assert(cache.GetMissRate() is: 1/2) // 1 miss / 2 gets
		Assert(result is: 'defaultValue')

		// user2 gets default value. it is not cached for user2
		result = UserSettingsCached().Get(.key, 'defaultValue', .user2)
		Assert(cache.GetMissRate() is: 2/3) // 2 miss / 3 gets
		Assert(result is: 'defaultValue')

		// user2 gets default value again. This time it is cached
		result = UserSettingsCached().Get(.key, 'defaultValue', .user2)
		Assert(cache.GetMissRate() is: 2/4) // 2 miss / 4 gets
		Assert(result is: 'defaultValue')

		// user1 saves a new value. This will reset the cache
		UserSettingsCached().Put(.key, 'newValue', .user1)
		Assert(cache.GetMissRate() is: 0)
		result = UserSettingsCached().Get(.key, 'defaultValue', .user1)
		Assert(result is: 'newValue')
		Assert(cache.GetMissRate() is: 1/1) // cache is reset. Back to 1 miss / 1 get

		// user1 gets the newValue again. This time it is cached
		result = UserSettingsCached().Get(.key, 'defaultValue', .user1)
		Assert(cache.GetMissRate() is: 1/2) // 1 miss / 2 gets
		Assert(result is: 'newValue')

		// user1 gets the value for a new key (not cached)
		result = UserSettingsCached().Get(.key2, 'defaultValue2', .user1)
		Assert(cache.GetMissRate() is: 2/3) // 2 miss / 3 gets
		Assert(result is: 'defaultValue2')

		// user2 gets a value cached prior to the reset. It should use the new default
		result = UserSettingsCached().Get(.key, 'defaultValueChanged', .user2)
		Assert(cache.GetMissRate() is: 3/4) // 2 miss / 4 gets
		Assert(result is: 'defaultValueChanged')


		// user2 gets the default value for the same key that user1 has cached
		result = UserSettingsCached().Get(.key2, 'defaultValue_user2', .user2)
		Assert(cache.GetMissRate() is: 4/5) // 2 miss / 3 gets
		Assert(result is: 'defaultValue_user2')

		// ensure Singleton.ResetAll will clear the cache
		Singleton.ResetAll()
		Assert(cache.GetMissRate() is: 0)
		}

	Test_CacheChanged()
		{
		UserSettingsCached.Reset()

		// test put before get
		UserSettingsCached().Put(.key3, 'initialValue', .user1)
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: 'initialValue')

		// test modifying and saving the cached value
		UserSettingsCached().Put(.key3, Object(), .user1)
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: #())

		// overwrite the result with a different object. UserSettingsCached should clear
		// the cache and the next put should get the changed value
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result.Readonly?())
		result = result.Copy()
		result.changedVal = 'abc'
		UserSettingsCached().Put(.key3, result, .user1)

		// same as above, without a put/save. Ensure that value that is cached can't
		// be changed by changing the return value
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: #(changedVal: 'abc'))
		result = result.Copy()
		result.changedVal = 'def'
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: #(changedVal: 'abc'))

		// test with multiple object levels
		result = Object(changedVal: Object(innerVal: '123'))
		UserSettingsCached().Put(.key3, result, .user1)
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: #(changedVal: #(innerVal: '123')))
		Assert(result.changedVal.Readonly?())

		// test overwriting the return value with a different type
		result = #20200101
		UserSettingsCached().Put(.key3, result, .user1)
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: #20200101)

		// test overwrite without put on non-object
		result = #20230101
		result = UserSettingsCached().Get(.key3, 'defaultValue', .user1)
		Assert(result is: #20200101)
		}

	Teardown()
		{
		UserSettingsCached.Reset()
		UserSettings.RemoveAllUsers(.key)
		UserSettings.RemoveAllUsers(.key2)
		UserSettings.RemoveAllUsers(.key3)
		if .created
			Database('destroy ' $ UserSettings.Table)
		}
	}
