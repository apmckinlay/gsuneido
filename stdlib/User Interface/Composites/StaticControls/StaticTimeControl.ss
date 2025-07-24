// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
StaticTextControl
	{
	Set(date)
		{
		super.Set(Date?(date) ? date.Time() : date)
		}
	}
