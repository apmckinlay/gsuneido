// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(x)
		{
		return String?(x) and x.Number?() and not x.Has?('e') and
			(x is '0' or not x.Prefix?('0'))
			? .convert(x) : x
		}
	convert(value)
		{
		newVal = value
		try
			newVal = Number(value)
		catch
			return value

		return value is String(newVal) ? newVal : value
		}
	}