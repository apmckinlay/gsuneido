// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return value <= args
		}
	Expected(args)
		{
		return "a value less than or equal to " $ Display(args)
		}
	}