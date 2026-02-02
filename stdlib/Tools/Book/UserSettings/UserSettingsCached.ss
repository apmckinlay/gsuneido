// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		.getResults = Object()
		}

	cache: Memoize
		{
		Func(key, def = false, user = false)
			{
			result = UserSettings.Get(key, def, user)
			if Object?(result)
				result.Set_readonly()
			return result
			}
		}

	Get(key, def = false, user = false)
		{
		// want to do this here to ensure .Func method signature matches with user = false
		// vs explicitly passing in Suneido.User
		if user is false
			user = Suneido.User

		result = (.cache)(key, def, user)
		.getResults[key] = result
		return result
		}

	Put(key, value, user = false, resetServer? = false)
		{
		if ((.getResults.Member?(key) and .getResults[key] isnt value) or
			not .getResults.Member?(key))
			{
			UserSettings.Put(key, value, user)
			(.cache).ResetCache()
			if resetServer?
				ServerEval('UserSettingsCached.Reset')
			}
		}

	Reset()
		{
		(.cache).ResetCache()
		super.Reset()
		}

	Remove(key)
		{
		UserSettings.Remove(key)
		(.cache).ResetCache()
		}
	}
