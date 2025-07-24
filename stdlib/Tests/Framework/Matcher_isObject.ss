// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args/*unused*/)
		{
		return Object?(value)
		}
	Expected(args/*unused*/)
		{
		return "an object"
		}
	}