// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args/*unused*/)
		{
		return Date?(value)
		}
	Expected(args/*unused*/)
		{
		return "a date"
		}
	}
