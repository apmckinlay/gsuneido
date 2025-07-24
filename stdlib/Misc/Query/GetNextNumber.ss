// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// WARNING: this is the old function, anything new should use GetNextNum instead
class
	{
	CallClass(table, field, log = false)
		{
		nextnum = false
		RetryTransaction()
			{|t|
			x = .Get(t, table)
			nextnum = x[field]++
			x.Update()
			}
		if log
			SuneidoLog("GetNextNumber from " $ table $ ": " $ String(nextnum))
		return nextnum
		}

	Get(t, table)
		{
		x = t.Query1(table)
		if x is false
			throw "GetNextNumber failed, no records in: " $ table
		return x
		}
	}