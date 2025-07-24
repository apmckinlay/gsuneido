// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Actual(value, args/*unused*/)
		{
		return "was " $ .DisplayValue(value)
		}

	DisplayValue(value)
		{
		displayVal = Display(value)
		if String?(value)
			displayVal = '\n' $ displayVal
		return displayVal.Ellipsis(1000) /*= max length */
		}
	}
