// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args/*unused*/)
		{
		return Number?(value)
		}
	Expected(args/*unused*/)
		{
		return "a number"
		}
	}