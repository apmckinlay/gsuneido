// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Table: 'user_settings'
	Get(key, def = false, user = false)
		{
		if .keyTooLong?(key)
			return def
		if user is false
			user = Suneido.User
		.Ensure()
		Transaction(read:)
			{|t|
			x = t.Query1(.query(key, user))
			if x is false
				x = t.Query1(.query(key, ''))
			}
		return x is false ? def : x.value
		}

	Put(key, value, user = false)
		{
		if .keyTooLong?(key)
			return
		if user is false
			user = Suneido.User
		.Ensure()
		query = .query(key, user)
		RetryTransaction()
			{ |t|
			if false is x = t.Query1(query)
				t.QueryOutput(.Table, Record(:user, :key, :value))
			else
				{
				x.value = value
				x.Update()
				}
			}
		.callObservers(key, value, user)
		}

	callObservers(key, value, user)
		{
		if Suneido.Member?("UserSettings_observers")
			for observer in Suneido.UserSettings_observers[key]
				observer(:value, :user)
		}

	AddObserver(key, observer)
		{
		if not Suneido.Member?("UserSettings_observers")
			Suneido.UserSettings_observers = Object().Set_default(Object())
		Suneido.UserSettings_observers[key].Add(observer)
		}

	Remove(key)
		{
		if .keyTooLong?(key)
			return
		.Ensure()
		QueryDo('delete ' $ .query(key))
		}
	RemoveAllUsers(key)
		{
		if .keyTooLong?(key)
			return
		.Ensure()
		QueryDo('delete ' $ .queryAllUsers(key))
		}
	keyTooLong?(key)
		{
		if key.Size() > 100 /*= key too large*/
			{
			SuneidoLog("ERROR: UserSettings key too large", params: key, calls:)
			return true
			}
		return false
		}
	Ensure()
		{
		Database('ensure ' $ .Table $ ' (user, key, value, usersetting_TS)
			index (key)
			key(user, key)')
		}
	query(key, user = false)
		{
		if user is false
			user = Suneido.User
		return .Table $ ' where
			user is ' $ Display(user) $ ' and
			key is ' $ Display(key)
		}
	queryAllUsers(key)
		{
		return .Table $ ' where user > "" and key is ' $ Display(key)
		}
	}