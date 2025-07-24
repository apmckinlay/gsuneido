// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// useful to determine virtual key values
// run it and then press the key combination
FieldControl
	{
	New()
		{ .SubClass() }
	KEYDOWN(wParam, lParam)
		{
		if wParam isnt 0x10 and wParam isnt 0x11
			Print('KEYDOWN', wParam.Hex())
		return 0
		}
	}