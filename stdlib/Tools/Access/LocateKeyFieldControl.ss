// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
KeyFieldControl
	{
	ComponentName: 'LocateKeyField'
	LBUTTONDOWN()
		{
		if GetFocus() is .Hwnd
			return 'callsuper'
		.SetFocus()
		.SelectAll()
		return 0
		}
	}