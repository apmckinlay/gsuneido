// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return args[0] <= value and value <= args[1]
		}
	Expected(args)
		{
		return "a value between " $ Display(args[0]) $ " and " $ Display(args[1])
		}
	}