// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
StaticTextControl
	{
	Set(date)
		{
		super.Set(Date?(date) ? date.LongDateTime() : date)
		}
	}
