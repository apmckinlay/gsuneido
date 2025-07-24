// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
OptionalNumberFormat
	{
	Print(x, y, w, h, data = '')
		{
		if data is 0
			data = ''
		super.Print(x, y, w, h, data)
		}
	}