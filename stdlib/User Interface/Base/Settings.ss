// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Get(setting)
		{
		settings = Suneido.GetDefault('Settings', [])
		return settings[setting]
		}

	Set(setting, value)
		{
		if not Suneido.Member?('Settings')
			Suneido.Settings = []
		Suneido.Settings[setting] = value
		}
	}
