// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
GetNextNum
	{
	Reserve(table, field = "nextnum")
		{
		return String(super.Reserve(table, field))
		}
	NumQuery(table, field, num)
		{
		val = String?(num) and num.Number?() ? Number(num) : num
		return super.NumQuery(table, field, val)
		}
	PutBack(num, table, field = "nextnum")
		{
		return super.PutBack(Number(num), table, field)
		}
	}