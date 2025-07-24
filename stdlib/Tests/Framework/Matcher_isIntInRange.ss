// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return Number?(value) and value.Int?() and
			args[0] <= value and value < args[1]
		}
	Expected(args)
		{
		return "an integer in the range [" $ args[0] $ ", " $ args[1] $ ")"
		}
	}