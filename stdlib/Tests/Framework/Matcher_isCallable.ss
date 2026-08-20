// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
MatcherWas
	{
	Match(value, args/*unused*/)
		{
		return Function?(value)
		}
	Expected(args/*unused*/)
		{
		return "a callable"
		}
	}
