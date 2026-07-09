// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(table, parent)
		{
		return .Get(table, parent)
		}

	Get(table, parent, t = false)
		{
		rec = [:parent]
		path = ''
		DoWithTran(t)
			{|t|
			while rec.parent isnt 0 and rec.parent isnt ''
				{
				if false is rec = t.Query1(table, num: rec.parent)
					break
				path = rec.name $ '/' $ path
				}
			}
		return path[.. -1]
		}
	}
