// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (table, t = false)
	{
	num = 1
	DoWithTran(t)
		{ |t| num = 1 + t.QueryMax(table, 'num', 0) }
	return num
	}