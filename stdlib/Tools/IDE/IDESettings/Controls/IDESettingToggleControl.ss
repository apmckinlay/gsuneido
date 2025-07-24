// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(name, defaultVal = false)
		{
		.Name = name
		value = IDESettings.Get(.Name, defaultVal) is false ? #Disable : #Enable
		.FindControl(#RadioButtons).Set(value)
		}

	Controls()
		{ return #(RadioButtons, Enable, Disable, horz:) }

	Get()
		{ return .FindControl(#RadioButtons).Get() is #Enable }
	}
