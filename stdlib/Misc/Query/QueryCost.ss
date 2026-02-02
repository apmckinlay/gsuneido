// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	OkForResetAll?: true
	Func(query)
		{
		WithQuery(query, { s = it.Strategy() })
		x = s.ExtractAll(`nrecs~ ([_0-9]+) cost~ ([_0-9]+)`)
		return Object(nrecs: Number(x[1]), cost: Number(x[2]))
		}
	}