// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args, _debugTest = false)
		{
		if args is value
			return true
		if debugTest and .diff?(value) and .diff?(args)
			Diff2Control('Assert FAILED', value, args, 'Actual', 'Expected')
		return false
		}
	diff?(x)
		{
		return String?(x) and x.Has?('\n')
		}
	Expected(args)
		{
		return .DisplayValue(args)
		}
	}