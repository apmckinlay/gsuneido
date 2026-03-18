// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// WARNING: this is the old function, anything new should use GetNextNum instead
class
	{
	CallClass(table, field)
		{
		nextnum = false
		RetryTransaction()
			{|t|
			x = .Get(t, table)
			nextnum = x[field]++
			x.Update()
			}
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