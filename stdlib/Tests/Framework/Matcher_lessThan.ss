// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args)
		{
		return value < args
		}
	Expected(args)
		{
		return "a value less than " $ Display(args)
		}
	}