// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Prompt: 'Date'
	Control: (ChooseDate)
	Format: (ShortDate)
	Encode(val, fmt = false)
		{
		if Date?(val)
			return val

		fmt = fmt is false ? Settings.Get('ShortDateFormat') : fmt
		// need 'try' to catch bad date literal exception,
		// this can happen on string values that start with '#'
		try
			x = val.Prefix?('#') ? Date(val) : Date(val, fmt)
		catch
			x = false
		return x isnt false ? x : val
		}
	}