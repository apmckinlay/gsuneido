// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		os = .getOsName()
		if .js?() or
			os.Has?('Windows 10') or os.Has?('Windows 11') or
			os is '' // when not able to get os name
			sid = 0
		else
			sid = .wts_GetSessionId()

		sid = sid is 0
			? .databaseSessionId() // ip address
			: 'wts' $ sid
		return sid
		}

	js?()
		{
		return Sys.SuneidoJs?()
		}

	getOsName()
		{
		return SystemInfo().OSName
		}

	wts_GetSessionId()
		{
		return WTS_GetSessionId()
		}

	databaseSessionId()
		{
		return Database.SessionId()
		}
	}
