// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	OkForResetAll?: true
	Func(query)
		{
		try
			Cursor(query) {|c| s = c.Strategy() }
		catch (unused, "invalid query")
			return Object(nrecs: 0, cost: .impossible)
		x = s.ExtractAll(`nrecs~ (\d+) cost~ (\d+)`)
		return Object(nrecs: Number(x[1]), cost: Number(x[2]))
		}
	impossible: 999999999999
	}