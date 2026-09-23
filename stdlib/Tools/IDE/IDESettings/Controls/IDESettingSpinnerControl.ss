// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(.Name, .defaultVal = 10)
		{
		}

	Controls()
		{
		if false is set = IDESettings.Get(.Name, .defaultVal)
			set = ""
		return [#Spinner, :set]
		}

	Get()
		{
		value = .FindControl(#Spinner).Get()
		return value is "" ? false : value
		}
	}
