// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args/*unused*/)
		{
		return Number?(value) and value.Int?()
		}
	Expected(args/*unused*/)
		{
		return "an integer"
		}
	}