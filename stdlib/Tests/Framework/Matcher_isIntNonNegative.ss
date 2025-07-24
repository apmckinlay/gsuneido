// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args/*unused*/)
		{
		return Number?(value) and value.Int?() and 0 <= value
		}
	Expected(args/*unused*/)
		{
		return "a non-negative integer"
		}
	}