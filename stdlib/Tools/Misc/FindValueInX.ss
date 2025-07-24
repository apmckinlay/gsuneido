// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Reference?(x, value)
		{
		return x is value or
			(String?(x) and x.Has?(',') and .checkCommaList(x, value)) or
			(Object?(x) and .traceObject(x, value))
		}

	checkCommaList(x, value)
		{
		for xVal in x.Split(',')
			if xVal.Trim() is value
				return true
		return false
		}

	traceObject(ob, value)
		{
		for x in ob // values
			if .Reference?(x, value)
				return true
		for x in ob.Members(named:)
			if .Reference?(x, value)
				return true
		return false
		}
	}